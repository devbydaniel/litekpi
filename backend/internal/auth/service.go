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

	"github.com/devbydaniel/litekpi/internal/platform/config"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrEmailNotVerified   = errors.New("email not verified")
	ErrUserNotFound       = errors.New("user not found")
	ErrTokenExpired       = errors.New("token has expired")
	ErrTokenUsed          = errors.New("token has already been used")
	ErrOAuthAccountNotFound = errors.New("oauth account not found")
)

const (
	bcryptCost               = 12
	verificationTokenExpiry  = 24 * time.Hour
	passwordResetTokenExpiry = 1 * time.Hour
	tokenLength              = 32
)

// Service handles authentication business logic.
type Service struct {
	repo        *Repository
	jwt         *JWTService
	email       *AuthEmailer
	googleOAuth *oauth2.Config
	githubOAuth *oauth2.Config
	appURL      string
}

// NewService creates a new auth service.
func NewService(repo *Repository, jwt *JWTService, email *AuthEmailer, cfg *config.Config) *Service {
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

// Register creates a new user with email, password, name, and organization.
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

	// Create organization and user in a transaction
	user, err := s.repo.CreateUserWithOrg(ctx, strings.ToLower(req.Email), req.Name, &passwordHash, req.OrganizationName)
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
	token, err := s.jwt.GenerateToken(user.ID, user.Email, user.OrganizationID, user.Role)
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

// OAuthResult represents the result of an OAuth callback.
type OAuthResult struct {
	// If user exists, AuthResponse is set
	AuthResponse *AuthResponse
	// If user needs setup, PendingSetup is set
	PendingSetup *OAuthPendingSetupResponse
}

// HandleGoogleCallback processes the Google OAuth callback.
func (s *Service) HandleGoogleCallback(ctx context.Context, code string) (*OAuthResult, error) {
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
func (s *Service) HandleGithubCallback(ctx context.Context, code string) (*OAuthResult, error) {
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

func (s *Service) handleOAuthUser(ctx context.Context, provider string, userInfo *OAuthUserInfo) (*OAuthResult, error) {
	// Check if OAuth account exists
	oauthAccount, err := s.repo.GetOAuthAccount(ctx, provider, userInfo.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get oauth account: %w", err)
	}

	// If OAuth account exists and has a linked user, log them in
	if oauthAccount != nil && oauthAccount.UserID != nil {
		user, err := s.repo.GetUserByID(ctx, *oauthAccount.UserID)
		if err != nil {
			return nil, fmt.Errorf("failed to get user: %w", err)
		}
		if user == nil {
			return nil, ErrUserNotFound
		}

		// Generate JWT token
		jwtToken, err := s.jwt.GenerateToken(user.ID, user.Email, user.OrganizationID, user.Role)
		if err != nil {
			return nil, fmt.Errorf("failed to generate token: %w", err)
		}

		return &OAuthResult{
			AuthResponse: &AuthResponse{
				User:  *user,
				Token: jwtToken,
			},
		}, nil
	}

	// Check if user with this email exists (link OAuth to existing user)
	existingUser, err := s.repo.GetUserByEmail(ctx, strings.ToLower(userInfo.Email))
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	if existingUser != nil {
		// Link OAuth account to existing user
		newOAuthAccount := &OAuthAccount{
			ID:             uuid.New(),
			UserID:         &existingUser.ID,
			Provider:       provider,
			ProviderUserID: userInfo.ID,
			ProviderEmail:  userInfo.Email,
			ProviderName:   userInfo.Name,
			CreatedAt:      time.Now(),
		}
		if err := s.repo.CreateOAuthAccount(ctx, newOAuthAccount); err != nil {
			return nil, fmt.Errorf("failed to create oauth account: %w", err)
		}

		// Mark email as verified (OAuth provider verified it)
		if !existingUser.EmailVerified {
			if err := s.repo.UpdateUserEmailVerified(ctx, existingUser.ID, true); err != nil {
				return nil, fmt.Errorf("failed to verify email: %w", err)
			}
			existingUser.EmailVerified = true
		}

		// Generate JWT token
		jwtToken, err := s.jwt.GenerateToken(existingUser.ID, existingUser.Email, existingUser.OrganizationID, existingUser.Role)
		if err != nil {
			return nil, fmt.Errorf("failed to generate token: %w", err)
		}

		return &OAuthResult{
			AuthResponse: &AuthResponse{
				User:  *existingUser,
				Token: jwtToken,
			},
		}, nil
	}

	// New OAuth user - create pending OAuth account and require setup
	pendingAccount, err := s.repo.CreatePendingOAuthAccount(ctx, provider, userInfo.ID, userInfo.Email, userInfo.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to create pending oauth account: %w", err)
	}

	// Generate short-lived setup token
	setupToken, err := s.jwt.GenerateOAuthSetupToken(pendingAccount.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate setup token: %w", err)
	}

	return &OAuthResult{
		PendingSetup: &OAuthPendingSetupResponse{
			PendingSetup: true,
			Token:        setupToken,
			Email:        userInfo.Email,
			ProviderName: userInfo.Name,
		},
	}, nil
}

// CompleteOAuthSetup finishes OAuth registration by creating org and user.
func (s *Service) CompleteOAuthSetup(ctx context.Context, req CompleteOAuthSetupRequest) (*AuthResponse, error) {
	// Validate setup token
	oauthAccountID, err := s.jwt.ValidateOAuthSetupToken(req.Token)
	if err != nil {
		return nil, ErrInvalidToken
	}

	// Get pending OAuth account
	oauthAccount, err := s.repo.GetOAuthAccountByID(ctx, oauthAccountID)
	if err != nil {
		return nil, fmt.Errorf("failed to get oauth account: %w", err)
	}
	if oauthAccount == nil {
		return nil, ErrOAuthAccountNotFound
	}

	// Check if already linked to a user
	if oauthAccount.UserID != nil {
		return nil, errors.New("oauth account already linked to a user")
	}

	// Check if email already exists
	existingUser, err := s.repo.GetUserByEmail(ctx, strings.ToLower(oauthAccount.ProviderEmail))
	if err != nil {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}
	if existingUser != nil {
		return nil, ErrEmailAlreadyExists
	}

	// Create organization and user (no password for OAuth users)
	user, err := s.repo.CreateUserWithOrg(ctx, strings.ToLower(oauthAccount.ProviderEmail), req.Name, nil, req.OrganizationName)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Mark email as verified (OAuth provider verified it)
	if err := s.repo.UpdateUserEmailVerified(ctx, user.ID, true); err != nil {
		return nil, fmt.Errorf("failed to verify email: %w", err)
	}
	user.EmailVerified = true

	// Link OAuth account to user
	if err := s.repo.LinkOAuthAccountToUser(ctx, oauthAccount.ID, user.ID); err != nil {
		return nil, fmt.Errorf("failed to link oauth account: %w", err)
	}

	// Generate JWT token
	jwtToken, err := s.jwt.GenerateToken(user.ID, user.Email, user.OrganizationID, user.Role)
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

// Invite expiry duration
const inviteExpiry = 7 * 24 * time.Hour

var (
	ErrInviteNotFound     = errors.New("invite not found")
	ErrInviteExpired      = errors.New("invite has expired")
	ErrInviteAlreadyUsed  = errors.New("invite has already been used")
	ErrLastAdmin          = errors.New("cannot remove or demote the last admin")
	ErrCannotRemoveSelf   = errors.New("cannot remove yourself")
	ErrUserAlreadyExists  = errors.New("user with this email already exists")
	ErrPendingInviteExists = errors.New("a pending invite already exists for this email")
)

// IsEmailEnabled returns whether email is configured.
func (s *Service) IsEmailEnabled() bool {
	return s.email.IsEnabled()
}

// CreateInvite creates a new invite.
func (s *Service) CreateInvite(ctx context.Context, req CreateInviteRequest, inviter *User) (*CreateInviteResponse, error) {
	email := strings.ToLower(req.Email)

	// Check if user already exists
	existingUser, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}
	if existingUser != nil {
		return nil, ErrUserAlreadyExists
	}

	// Check if there's already a pending invite
	existingInvite, err := s.repo.GetPendingInviteByEmail(ctx, inviter.OrganizationID, email)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing invite: %w", err)
	}
	if existingInvite != nil {
		return nil, ErrPendingInviteExists
	}

	// Generate token
	token, err := generateSecureToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Create invite
	expiresAt := time.Now().Add(inviteExpiry)
	invite, err := s.repo.CreateInvite(ctx, inviter.OrganizationID, email, req.Role, token, inviter.ID, expiresAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create invite: %w", err)
	}

	// Try to send email
	if s.email.IsEnabled() {
		org, err := s.repo.GetOrganizationByID(ctx, inviter.OrganizationID)
		if err != nil {
			return nil, fmt.Errorf("failed to get organization: %w", err)
		}
		if err := s.email.SendInviteEmail(email, token, inviter.Name, org.Name); err != nil {
			// Log but don't fail if email fails
			fmt.Printf("failed to send invite email: %v\n", err)
		}
		return &CreateInviteResponse{Invite: *invite}, nil
	}

	// Email not configured - return invite URL
	inviteURL := fmt.Sprintf("%s/accept-invite?token=%s", s.appURL, token)
	return &CreateInviteResponse{Invite: *invite, InviteURL: &inviteURL}, nil
}

// ListInvites lists pending invites for an organization.
func (s *Service) ListInvites(ctx context.Context, orgID uuid.UUID) ([]InviteWithInviter, error) {
	invites, err := s.repo.ListPendingInvites(ctx, orgID)
	if err != nil {
		return nil, fmt.Errorf("failed to list invites: %w", err)
	}
	if invites == nil {
		return []InviteWithInviter{}, nil
	}
	return invites, nil
}

// CancelInvite cancels a pending invite.
func (s *Service) CancelInvite(ctx context.Context, inviteID uuid.UUID, orgID uuid.UUID) error {
	invite, err := s.repo.GetInviteByID(ctx, inviteID)
	if err != nil {
		return fmt.Errorf("failed to get invite: %w", err)
	}
	if invite == nil {
		return ErrInviteNotFound
	}
	if invite.OrganizationID != orgID {
		return ErrInviteNotFound
	}
	if invite.AcceptedAt != nil {
		return ErrInviteAlreadyUsed
	}

	if err := s.repo.DeleteInvite(ctx, inviteID); err != nil {
		return fmt.Errorf("failed to delete invite: %w", err)
	}
	return nil
}

// ValidateInvite validates an invite token and returns info.
func (s *Service) ValidateInvite(ctx context.Context, token string) (*ValidateInviteResponse, error) {
	invite, err := s.repo.GetInviteByToken(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("failed to get invite: %w", err)
	}
	if invite == nil {
		return &ValidateInviteResponse{Valid: false}, nil
	}
	if invite.AcceptedAt != nil {
		return &ValidateInviteResponse{Valid: false}, nil
	}
	if time.Now().After(invite.ExpiresAt) {
		return &ValidateInviteResponse{Valid: false}, nil
	}

	// Get organization and inviter info
	org, err := s.repo.GetOrganizationByID(ctx, invite.OrganizationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get organization: %w", err)
	}
	inviter, err := s.repo.GetUserByID(ctx, invite.InvitedBy)
	if err != nil {
		return nil, fmt.Errorf("failed to get inviter: %w", err)
	}

	return &ValidateInviteResponse{
		Valid:            true,
		Email:            invite.Email,
		OrganizationName: org.Name,
		Role:             invite.Role,
		InviterName:      inviter.Name,
	}, nil
}

// AcceptInvite accepts an invite and creates the user.
func (s *Service) AcceptInvite(ctx context.Context, req AcceptInviteRequest) (*User, error) {
	invite, err := s.repo.GetInviteByToken(ctx, req.Token)
	if err != nil {
		return nil, fmt.Errorf("failed to get invite: %w", err)
	}
	if invite == nil {
		return nil, ErrInviteNotFound
	}
	if invite.AcceptedAt != nil {
		return nil, ErrInviteAlreadyUsed
	}
	if time.Now().After(invite.ExpiresAt) {
		return nil, ErrInviteExpired
	}

	// Check if email already exists (in case someone registered in the meantime)
	existingUser, err := s.repo.GetUserByEmail(ctx, invite.Email)
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
	hash := string(hashedPassword)

	// Create user
	user, err := s.repo.CreateUser(ctx, invite.Email, req.Name, &hash, invite.OrganizationID, invite.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Mark email as verified (trusted from invite)
	if err := s.repo.UpdateUserEmailVerified(ctx, user.ID, true); err != nil {
		return nil, fmt.Errorf("failed to verify email: %w", err)
	}
	user.EmailVerified = true

	// Mark invite as accepted
	if err := s.repo.MarkInviteAccepted(ctx, invite.ID); err != nil {
		return nil, fmt.Errorf("failed to mark invite accepted: %w", err)
	}

	return user, nil
}

// ListUsers lists users in an organization.
func (s *Service) ListUsers(ctx context.Context, orgID uuid.UUID) ([]User, error) {
	users, err := s.repo.ListUsersByOrg(ctx, orgID)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	if users == nil {
		return []User{}, nil
	}
	return users, nil
}

// UpdateUserRole updates a user's role.
func (s *Service) UpdateUserRole(ctx context.Context, userID uuid.UUID, role Role, requestingUser *User) error {
	// Get target user
	targetUser, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	if targetUser == nil {
		return ErrUserNotFound
	}
	if targetUser.OrganizationID != requestingUser.OrganizationID {
		return ErrUserNotFound
	}

	// If demoting an admin, check if it's the last one
	if targetUser.Role == RoleAdmin && role != RoleAdmin {
		count, err := s.repo.CountAdmins(ctx, requestingUser.OrganizationID)
		if err != nil {
			return fmt.Errorf("failed to count admins: %w", err)
		}
		if count <= 1 {
			return ErrLastAdmin
		}
	}

	if err := s.repo.UpdateUserRole(ctx, userID, role); err != nil {
		return fmt.Errorf("failed to update user role: %w", err)
	}
	return nil
}

// RemoveUser removes a user from the organization.
func (s *Service) RemoveUser(ctx context.Context, userID uuid.UUID, requestingUser *User) error {
	// Cannot remove self
	if userID == requestingUser.ID {
		return ErrCannotRemoveSelf
	}

	// Get target user
	targetUser, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	if targetUser == nil {
		return ErrUserNotFound
	}
	if targetUser.OrganizationID != requestingUser.OrganizationID {
		return ErrUserNotFound
	}

	// If removing an admin, check if it's the last one
	if targetUser.Role == RoleAdmin {
		count, err := s.repo.CountAdmins(ctx, requestingUser.OrganizationID)
		if err != nil {
			return fmt.Errorf("failed to count admins: %w", err)
		}
		if count <= 1 {
			return ErrLastAdmin
		}
	}

	if err := s.repo.DeleteUser(ctx, userID); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}
