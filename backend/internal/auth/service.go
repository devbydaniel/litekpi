package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"

	"github.com/trackable/trackable/internal/platform/config"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrEmailNotVerified   = errors.New("email not verified")
	ErrUserNotFound       = errors.New("user not found")
	ErrTokenExpired       = errors.New("token has expired")
	ErrTokenUsed          = errors.New("token has already been used")
)

const (
	bcryptCost                     = 12
	verificationTokenExpiry        = 24 * time.Hour
	passwordResetTokenExpiry       = 1 * time.Hour
	tokenLength                    = 32
)

// Service handles authentication business logic.
type Service struct {
	repo        *Repository
	jwt         *JWTService
	email       *EmailService
	googleOAuth *oauth2.Config
	githubOAuth *oauth2.Config
	appURL      string
}

// NewService creates a new auth service.
func NewService(repo *Repository, jwt *JWTService, email *EmailService, cfg *config.Config) *Service {
	svc := &Service{
		repo:   repo,
		jwt:    jwt,
		email:  email,
		appURL: strings.TrimSuffix(cfg.AppURL, "/"),
	}

	// Configure Google OAuth if credentials provided
	if cfg.OAuth.GoogleClientID != "" && cfg.OAuth.GoogleClientSecret != "" {
		svc.googleOAuth = &oauth2.Config{
			ClientID:     cfg.OAuth.GoogleClientID,
			ClientSecret: cfg.OAuth.GoogleClientSecret,
			RedirectURL:  strings.TrimSuffix(cfg.APIURL, "/") + "/api/v1/auth/google/callback",
			Scopes:       []string{"email", "profile"},
			Endpoint:     google.Endpoint,
		}
	}

	// Configure GitHub OAuth if credentials provided
	if cfg.OAuth.GithubClientID != "" && cfg.OAuth.GithubClientSecret != "" {
		svc.githubOAuth = &oauth2.Config{
			ClientID:     cfg.OAuth.GithubClientID,
			ClientSecret: cfg.OAuth.GithubClientSecret,
			RedirectURL:  strings.TrimSuffix(cfg.APIURL, "/") + "/api/v1/auth/github/callback",
			Scopes:       []string{"user:email"},
			Endpoint:     github.Endpoint,
		}
	}

	return svc
}

// Register creates a new user with email and password.
func (s *Service) Register(ctx context.Context, req RegisterRequest) (*User, error) {
	// Check if email already exists
	existingUser, err := s.repo.GetUserByEmail(ctx, strings.ToLower(req.Email))
	if err != nil {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}
	if existingUser != nil {
		return nil, ErrEmailAlreadyExists
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcryptCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}
	passwordHash := string(hashedPassword)

	// Create user
	user, err := s.repo.CreateUser(ctx, strings.ToLower(req.Email), &passwordHash)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Send verification email
	if err := s.SendVerificationEmail(ctx, user.ID); err != nil {
		// Log error but don't fail registration
		fmt.Printf("failed to send verification email: %v\n", err)
	}

	return user, nil
}

// Login authenticates a user with email and password.
func (s *Service) Login(ctx context.Context, req LoginRequest) (*AuthResponse, error) {
	user, err := s.repo.GetUserByEmail(ctx, strings.ToLower(req.Email))
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return nil, ErrInvalidCredentials
	}

	// Check password
	if user.PasswordHash == nil {
		return nil, ErrInvalidCredentials // OAuth-only user
	}
	if err := bcrypt.CompareHashAndPassword([]byte(*user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	// Check email verification
	if !user.EmailVerified {
		return nil, ErrEmailNotVerified
	}

	// Generate token
	token, err := s.jwt.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &AuthResponse{
		User:  *user,
		Token: token,
	}, nil
}

// SendVerificationEmail sends a verification email to the user.
func (s *Service) SendVerificationEmail(ctx context.Context, userID uuid.UUID) error {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return ErrUserNotFound
	}

	// Delete existing tokens
	if err := s.repo.DeleteEmailVerificationTokensByUserID(ctx, userID); err != nil {
		return fmt.Errorf("failed to delete existing tokens: %w", err)
	}

	// Generate token
	token, err := generateSecureToken()
	if err != nil {
		return fmt.Errorf("failed to generate token: %w", err)
	}

	// Save token
	expiresAt := time.Now().Add(verificationTokenExpiry)
	if err := s.repo.CreateEmailVerificationToken(ctx, userID, token, expiresAt); err != nil {
		return fmt.Errorf("failed to save token: %w", err)
	}

	// Send email
	return s.email.SendVerificationEmail(user.Email, token)
}

// VerifyEmail verifies a user's email using a token.
func (s *Service) VerifyEmail(ctx context.Context, token string) error {
	evt, err := s.repo.GetEmailVerificationToken(ctx, token)
	if err != nil {
		return fmt.Errorf("failed to get token: %w", err)
	}
	if evt == nil {
		return ErrInvalidToken
	}

	// Check expiry
	if time.Now().After(evt.ExpiresAt) {
		return ErrTokenExpired
	}

	// Update user
	if err := s.repo.UpdateUserEmailVerified(ctx, evt.UserID, true); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	// Delete token
	if err := s.repo.DeleteEmailVerificationTokensByUserID(ctx, evt.UserID); err != nil {
		return fmt.Errorf("failed to delete tokens: %w", err)
	}

	return nil
}

// ResendVerificationEmail resends the verification email.
func (s *Service) ResendVerificationEmail(ctx context.Context, email string) error {
	user, err := s.repo.GetUserByEmail(ctx, strings.ToLower(email))
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return nil // Don't reveal if user exists
	}
	if user.EmailVerified {
		return nil // Already verified
	}

	return s.SendVerificationEmail(ctx, user.ID)
}

// RequestPasswordReset initiates a password reset.
func (s *Service) RequestPasswordReset(ctx context.Context, email string) error {
	user, err := s.repo.GetUserByEmail(ctx, strings.ToLower(email))
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return nil // Don't reveal if user exists
	}

	// Delete existing tokens
	if err := s.repo.DeletePasswordResetTokensByUserID(ctx, user.ID); err != nil {
		return fmt.Errorf("failed to delete existing tokens: %w", err)
	}

	// Generate token
	token, err := generateSecureToken()
	if err != nil {
		return fmt.Errorf("failed to generate token: %w", err)
	}

	// Save token
	expiresAt := time.Now().Add(passwordResetTokenExpiry)
	if err := s.repo.CreatePasswordResetToken(ctx, user.ID, token, expiresAt); err != nil {
		return fmt.Errorf("failed to save token: %w", err)
	}

	// Send email
	return s.email.SendPasswordResetEmail(user.Email, token)
}

// ResetPassword resets a user's password using a token.
func (s *Service) ResetPassword(ctx context.Context, token, newPassword string) error {
	prt, err := s.repo.GetPasswordResetToken(ctx, token)
	if err != nil {
		return fmt.Errorf("failed to get token: %w", err)
	}
	if prt == nil {
		return ErrInvalidToken
	}

	// Check if used
	if prt.Used {
		return ErrTokenUsed
	}

	// Check expiry
	if time.Now().After(prt.ExpiresAt) {
		return ErrTokenExpired
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcryptCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Update password
	if err := s.repo.UpdateUserPassword(ctx, prt.UserID, string(hashedPassword)); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	// Mark token as used
	if err := s.repo.MarkPasswordResetTokenUsed(ctx, prt.ID); err != nil {
		return fmt.Errorf("failed to mark token used: %w", err)
	}

	// Also verify email if not already verified
	if err := s.repo.UpdateUserEmailVerified(ctx, prt.UserID, true); err != nil {
		return fmt.Errorf("failed to verify email: %w", err)
	}

	return nil
}

// GetGoogleAuthURL returns the Google OAuth authorization URL.
func (s *Service) GetGoogleAuthURL(state string) string {
	if s.googleOAuth == nil {
		return ""
	}
	return s.googleOAuth.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

// GetGithubAuthURL returns the GitHub OAuth authorization URL.
func (s *Service) GetGithubAuthURL(state string) string {
	if s.githubOAuth == nil {
		return ""
	}
	return s.githubOAuth.AuthCodeURL(state)
}

// IsGoogleOAuthEnabled returns true if Google OAuth is configured.
func (s *Service) IsGoogleOAuthEnabled() bool {
	return s.googleOAuth != nil
}

// IsGithubOAuthEnabled returns true if GitHub OAuth is configured.
func (s *Service) IsGithubOAuthEnabled() bool {
	return s.githubOAuth != nil
}

// HandleGoogleCallback processes the Google OAuth callback.
func (s *Service) HandleGoogleCallback(ctx context.Context, code string) (*AuthResponse, error) {
	if s.googleOAuth == nil {
		return nil, errors.New("google oauth not configured")
	}

	// Exchange code for token
	token, err := s.googleOAuth.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}

	// Get user info
	userInfo, err := s.getGoogleUserInfo(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	return s.handleOAuthUser(ctx, "google", userInfo)
}

// HandleGithubCallback processes the GitHub OAuth callback.
func (s *Service) HandleGithubCallback(ctx context.Context, code string) (*AuthResponse, error) {
	if s.githubOAuth == nil {
		return nil, errors.New("github oauth not configured")
	}

	// Exchange code for token
	token, err := s.githubOAuth.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}

	// Get user info
	userInfo, err := s.getGithubUserInfo(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	return s.handleOAuthUser(ctx, "github", userInfo)
}

func (s *Service) getGoogleUserInfo(ctx context.Context, token *oauth2.Token) (*OAuthUserInfo, error) {
	client := s.googleOAuth.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data struct {
		ID    string `json:"id"`
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	return &OAuthUserInfo{
		ID:    data.ID,
		Email: data.Email,
		Name:  data.Name,
	}, nil
}

func (s *Service) getGithubUserInfo(ctx context.Context, token *oauth2.Token) (*OAuthUserInfo, error) {
	client := s.githubOAuth.Client(ctx, token)

	// Get user profile
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var userData struct {
		ID    int    `json:"id"`
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	if err := json.Unmarshal(body, &userData); err != nil {
		return nil, err
	}

	email := userData.Email

	// If email is empty, fetch from emails endpoint
	if email == "" {
		emailResp, err := client.Get("https://api.github.com/user/emails")
		if err != nil {
			return nil, err
		}
		defer emailResp.Body.Close()

		emailBody, err := io.ReadAll(emailResp.Body)
		if err != nil {
			return nil, err
		}

		var emails []struct {
			Email    string `json:"email"`
			Primary  bool   `json:"primary"`
			Verified bool   `json:"verified"`
		}
		if err := json.Unmarshal(emailBody, &emails); err != nil {
			return nil, err
		}

		// Find primary verified email
		for _, e := range emails {
			if e.Primary && e.Verified {
				email = e.Email
				break
			}
		}
	}

	return &OAuthUserInfo{
		ID:    fmt.Sprintf("%d", userData.ID),
		Email: email,
		Name:  userData.Name,
	}, nil
}

func (s *Service) handleOAuthUser(ctx context.Context, provider string, userInfo *OAuthUserInfo) (*AuthResponse, error) {
	// Check if OAuth account exists
	oauthAccount, err := s.repo.GetOAuthAccount(ctx, provider, userInfo.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get oauth account: %w", err)
	}

	var user *User

	if oauthAccount != nil {
		// Get existing user
		user, err = s.repo.GetUserByID(ctx, oauthAccount.UserID)
		if err != nil {
			return nil, fmt.Errorf("failed to get user: %w", err)
		}
	} else {
		// Check if user with this email exists
		user, err = s.repo.GetUserByEmail(ctx, strings.ToLower(userInfo.Email))
		if err != nil {
			return nil, fmt.Errorf("failed to get user by email: %w", err)
		}

		if user == nil {
			// Create new user (no password, OAuth only)
			user, err = s.repo.CreateUser(ctx, strings.ToLower(userInfo.Email), nil)
			if err != nil {
				return nil, fmt.Errorf("failed to create user: %w", err)
			}
		}

		// Mark email as verified (OAuth provider verified it)
		if !user.EmailVerified {
			if err := s.repo.UpdateUserEmailVerified(ctx, user.ID, true); err != nil {
				return nil, fmt.Errorf("failed to verify email: %w", err)
			}
			user.EmailVerified = true
		}

		// Create OAuth account link
		oauthAccount = &OAuthAccount{
			ID:             uuid.New(),
			UserID:         user.ID,
			Provider:       provider,
			ProviderUserID: userInfo.ID,
			ProviderEmail:  userInfo.Email,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}
		if err := s.repo.CreateOAuthAccount(ctx, oauthAccount); err != nil {
			return nil, fmt.Errorf("failed to create oauth account: %w", err)
		}
	}

	// Generate JWT token
	jwtToken, err := s.jwt.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &AuthResponse{
		User:  *user,
		Token: jwtToken,
	}, nil
}

// GetUserByID retrieves a user by ID.
func (s *Service) GetUserByID(ctx context.Context, id uuid.UUID) (*User, error) {
	return s.repo.GetUserByID(ctx, id)
}

// GetAppURL returns the configured app URL.
func (s *Service) GetAppURL() string {
	return s.appURL
}

func generateSecureToken() (string, error) {
	bytes := make([]byte, tokenLength)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// Middleware returns the auth middleware function.
func (s *Service) Middleware(next http.Handler) http.Handler {
	return AuthMiddleware(s.jwt, s.repo)(next)
}
