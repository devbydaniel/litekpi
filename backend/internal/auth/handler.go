package auth

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
)

// Handler handles HTTP requests for authentication.
type Handler struct {
	service *Service
}

// NewHandler creates a new auth handler.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Register handles user registration.
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validate input
	if req.Email == "" {
		respondError(w, http.StatusBadRequest, "email is required")
		return
	}
	if req.Password == "" {
		respondError(w, http.StatusBadRequest, "password is required")
		return
	}
	if len(req.Password) < 8 {
		respondError(w, http.StatusBadRequest, "password must be at least 8 characters")
		return
	}

	user, err := h.service.Register(r.Context(), req)
	if err != nil {
		if errors.Is(err, ErrEmailAlreadyExists) {
			respondError(w, http.StatusConflict, "email already exists")
			return
		}
		respondError(w, http.StatusInternalServerError, "failed to register user")
		return
	}

	respondJSON(w, http.StatusCreated, map[string]interface{}{
		"message": "Registration successful. Please check your email to verify your account.",
		"user":    user,
	})
}

// Login handles user login.
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validate input
	if req.Email == "" || req.Password == "" {
		respondError(w, http.StatusBadRequest, "email and password are required")
		return
	}

	resp, err := h.service.Login(r.Context(), req)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			respondError(w, http.StatusUnauthorized, "invalid email or password")
			return
		}
		if errors.Is(err, ErrEmailNotVerified) {
			respondError(w, http.StatusForbidden, "please verify your email before logging in")
			return
		}
		respondError(w, http.StatusInternalServerError, "failed to login")
		return
	}

	respondJSON(w, http.StatusOK, resp)
}

// VerifyEmail handles email verification.
func (h *Handler) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	var req VerifyEmailRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Token == "" {
		respondError(w, http.StatusBadRequest, "token is required")
		return
	}

	err := h.service.VerifyEmail(r.Context(), req.Token)
	if err != nil {
		if errors.Is(err, ErrInvalidToken) {
			respondError(w, http.StatusBadRequest, "invalid verification token")
			return
		}
		if errors.Is(err, ErrTokenExpired) {
			respondError(w, http.StatusBadRequest, "verification token has expired")
			return
		}
		respondError(w, http.StatusInternalServerError, "failed to verify email")
		return
	}

	respondJSON(w, http.StatusOK, MessageResponse{Message: "Email verified successfully"})
}

// ForgotPassword handles password reset requests.
func (h *Handler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req ForgotPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Email == "" {
		respondError(w, http.StatusBadRequest, "email is required")
		return
	}

	// Always return success to not reveal if email exists
	_ = h.service.RequestPasswordReset(r.Context(), req.Email)

	respondJSON(w, http.StatusOK, MessageResponse{Message: "If an account with that email exists, a password reset link has been sent"})
}

// ResetPassword handles password reset.
func (h *Handler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req ResetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Token == "" {
		respondError(w, http.StatusBadRequest, "token is required")
		return
	}
	if req.NewPassword == "" {
		respondError(w, http.StatusBadRequest, "new password is required")
		return
	}
	if len(req.NewPassword) < 8 {
		respondError(w, http.StatusBadRequest, "password must be at least 8 characters")
		return
	}

	err := h.service.ResetPassword(r.Context(), req.Token, req.NewPassword)
	if err != nil {
		if errors.Is(err, ErrInvalidToken) {
			respondError(w, http.StatusBadRequest, "invalid reset token")
			return
		}
		if errors.Is(err, ErrTokenExpired) {
			respondError(w, http.StatusBadRequest, "reset token has expired")
			return
		}
		if errors.Is(err, ErrTokenUsed) {
			respondError(w, http.StatusBadRequest, "reset token has already been used")
			return
		}
		respondError(w, http.StatusInternalServerError, "failed to reset password")
		return
	}

	respondJSON(w, http.StatusOK, MessageResponse{Message: "Password reset successfully"})
}

// ResendVerification handles resending verification emails.
func (h *Handler) ResendVerification(w http.ResponseWriter, r *http.Request) {
	var req ResendVerificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Email == "" {
		respondError(w, http.StatusBadRequest, "email is required")
		return
	}

	// Always return success to not reveal if email exists
	_ = h.service.ResendVerificationEmail(r.Context(), req.Email)

	respondJSON(w, http.StatusOK, MessageResponse{Message: "If an unverified account with that email exists, a verification link has been sent"})
}

// GoogleAuth initiates Google OAuth flow.
func (h *Handler) GoogleAuth(w http.ResponseWriter, r *http.Request) {
	if !h.service.IsGoogleOAuthEnabled() {
		respondError(w, http.StatusNotFound, "google oauth not configured")
		return
	}

	state := generateOAuthState()
	authURL := h.service.GetGoogleAuthURL(state)

	// Store state in cookie for validation
	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		Path:     "/",
		MaxAge:   600, // 10 minutes
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

// GoogleCallback handles Google OAuth callback.
func (h *Handler) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	// Validate state
	stateCookie, err := r.Cookie("oauth_state")
	if err != nil || stateCookie.Value != r.URL.Query().Get("state") {
		h.redirectWithError(w, r, "invalid oauth state")
		return
	}

	// Clear state cookie
	http.SetCookie(w, &http.Cookie{
		Name:   "oauth_state",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	code := r.URL.Query().Get("code")
	if code == "" {
		h.redirectWithError(w, r, "missing authorization code")
		return
	}

	resp, err := h.service.HandleGoogleCallback(r.Context(), code)
	if err != nil {
		h.redirectWithError(w, r, "authentication failed")
		return
	}

	h.redirectWithAuth(w, r, resp)
}

// GithubAuth initiates GitHub OAuth flow.
func (h *Handler) GithubAuth(w http.ResponseWriter, r *http.Request) {
	if !h.service.IsGithubOAuthEnabled() {
		respondError(w, http.StatusNotFound, "github oauth not configured")
		return
	}

	state := generateOAuthState()
	authURL := h.service.GetGithubAuthURL(state)

	// Store state in cookie for validation
	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		Path:     "/",
		MaxAge:   600, // 10 minutes
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

// GithubCallback handles GitHub OAuth callback.
func (h *Handler) GithubCallback(w http.ResponseWriter, r *http.Request) {
	// Validate state
	stateCookie, err := r.Cookie("oauth_state")
	if err != nil || stateCookie.Value != r.URL.Query().Get("state") {
		h.redirectWithError(w, r, "invalid oauth state")
		return
	}

	// Clear state cookie
	http.SetCookie(w, &http.Cookie{
		Name:   "oauth_state",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	code := r.URL.Query().Get("code")
	if code == "" {
		h.redirectWithError(w, r, "missing authorization code")
		return
	}

	resp, err := h.service.HandleGithubCallback(r.Context(), code)
	if err != nil {
		h.redirectWithError(w, r, "authentication failed")
		return
	}

	h.redirectWithAuth(w, r, resp)
}

// Me returns the current user.
func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	user := UserFromContext(r.Context())
	if user == nil {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	respondJSON(w, http.StatusOK, user)
}

// Logout handles user logout (client-side token invalidation).
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, MessageResponse{Message: "Logged out successfully"})
}

func (h *Handler) redirectWithAuth(w http.ResponseWriter, r *http.Request, resp *AuthResponse) {
	// Encode user data as base64 JSON
	userData, _ := json.Marshal(resp.User)
	userEncoded := base64.URLEncoding.EncodeToString(userData)

	redirectURL := h.service.GetAppURL() + "/auth/callback?token=" + url.QueryEscape(resp.Token) + "&user=" + url.QueryEscape(userEncoded)
	http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
}

func (h *Handler) redirectWithError(w http.ResponseWriter, r *http.Request, message string) {
	redirectURL := h.service.GetAppURL() + "/login?error=" + url.QueryEscape(message)
	http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
}

func generateOAuthState() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
