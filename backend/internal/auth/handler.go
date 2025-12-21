package auth

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
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
//
//	@Summary		Register a new user
//	@Description	Create a new user account with email, password, name, and organization
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		RegisterRequest		true	"Registration data"
//	@Success		201		{object}	RegisterResponse
//	@Failure		400		{object}	ErrorResponse
//	@Failure		409		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Router			/auth/register [post]
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
	if req.Name == "" {
		respondError(w, http.StatusBadRequest, "name is required")
		return
	}
	if req.OrganizationName == "" {
		respondError(w, http.StatusBadRequest, "organization name is required")
		return
	}

	user, err := h.service.Register(r.Context(), req)
	if err != nil {
		if errors.Is(err, ErrEmailAlreadyExists) {
			respondError(w, http.StatusConflict, "email already exists")
			return
		}
		log.Printf("registration error: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to register user")
		return
	}

	respondJSON(w, http.StatusCreated, map[string]interface{}{
		"message": "Registration successful. Please check your email to verify your account.",
		"user":    user,
	})
}

// Login handles user login.
//
//	@Summary		Login user
//	@Description	Authenticate user with email and password
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		LoginRequest	true	"Login credentials"
//	@Success		200		{object}	AuthResponse
//	@Failure		400		{object}	ErrorResponse
//	@Failure		401		{object}	ErrorResponse
//	@Failure		403		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Router			/auth/login [post]
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
//
//	@Summary		Verify email
//	@Description	Verify user email address with token
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		VerifyEmailRequest	true	"Verification token"
//	@Success		200		{object}	MessageResponse
//	@Failure		400		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Router			/auth/verify-email [post]
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
//
//	@Summary		Request password reset
//	@Description	Send password reset email to user
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		ForgotPasswordRequest	true	"User email"
//	@Success		200		{object}	MessageResponse
//	@Failure		400		{object}	ErrorResponse
//	@Router			/auth/forgot-password [post]
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
//
//	@Summary		Reset password
//	@Description	Reset user password with token
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		ResetPasswordRequest	true	"Reset token and new password"
//	@Success		200		{object}	MessageResponse
//	@Failure		400		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Router			/auth/reset-password [post]
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
//
//	@Summary		Resend verification email
//	@Description	Resend email verification link to user
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		ResendVerificationRequest	true	"User email"
//	@Success		200		{object}	MessageResponse
//	@Failure		400		{object}	ErrorResponse
//	@Router			/auth/resend-verification [post]
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

	result, err := h.service.HandleGoogleCallback(r.Context(), code)
	if err != nil {
		h.redirectWithError(w, r, "authentication failed")
		return
	}

	h.redirectWithOAuthResult(w, r, result)
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

	result, err := h.service.HandleGithubCallback(r.Context(), code)
	if err != nil {
		h.redirectWithError(w, r, "authentication failed")
		return
	}

	h.redirectWithOAuthResult(w, r, result)
}

// CompleteOAuthSetup handles completing OAuth registration with org/name.
//
//	@Summary		Complete OAuth setup
//	@Description	Complete OAuth registration by providing name and organization
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		CompleteOAuthSetupRequest	true	"Setup data"
//	@Success		200		{object}	AuthResponse
//	@Failure		400		{object}	ErrorResponse
//	@Failure		409		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Router			/auth/complete-oauth-setup [post]
func (h *Handler) CompleteOAuthSetup(w http.ResponseWriter, r *http.Request) {
	var req CompleteOAuthSetupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validate input
	if req.Token == "" {
		respondError(w, http.StatusBadRequest, "token is required")
		return
	}
	if req.Name == "" {
		respondError(w, http.StatusBadRequest, "name is required")
		return
	}
	if req.OrganizationName == "" {
		respondError(w, http.StatusBadRequest, "organization name is required")
		return
	}

	resp, err := h.service.CompleteOAuthSetup(r.Context(), req)
	if err != nil {
		if errors.Is(err, ErrInvalidToken) {
			respondError(w, http.StatusBadRequest, "invalid or expired setup token")
			return
		}
		if errors.Is(err, ErrOAuthAccountNotFound) {
			respondError(w, http.StatusBadRequest, "oauth account not found")
			return
		}
		if errors.Is(err, ErrEmailAlreadyExists) {
			respondError(w, http.StatusConflict, "email already exists")
			return
		}
		log.Printf("complete oauth setup error: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to complete setup")
		return
	}

	respondJSON(w, http.StatusOK, resp)
}

// Me returns the current user.
//
//	@Summary		Get current user
//	@Description	Get the currently authenticated user's profile
//	@Tags			auth
//	@Produce		json
//	@Success		200	{object}	User
//	@Failure		401	{object}	ErrorResponse
//	@Security		BearerAuth
//	@Router			/auth/me [get]
func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	user := UserFromContext(r.Context())
	if user == nil {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	respondJSON(w, http.StatusOK, user)
}

// Logout handles user logout (client-side token invalidation).
//
//	@Summary		Logout user
//	@Description	Logout the current user (client should discard token)
//	@Tags			auth
//	@Produce		json
//	@Success		200	{object}	MessageResponse
//	@Security		BearerAuth
//	@Router			/auth/logout [post]
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, MessageResponse{Message: "Logged out successfully"})
}

// GetEmailConfig returns email configuration status.
//
//	@Summary		Get email config status
//	@Description	Check if email is configured
//	@Tags			auth
//	@Produce		json
//	@Success		200	{object}	EmailConfigResponse
//	@Security		BearerAuth
//	@Router			/auth/email-config [get]
func (h *Handler) GetEmailConfig(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, EmailConfigResponse{Enabled: h.service.IsEmailEnabled()})
}

// CreateInvite creates a new invite.
//
//	@Summary		Create invite
//	@Description	Create a new user invite
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		CreateInviteRequest	true	"Invite data"
//	@Success		201		{object}	CreateInviteResponse
//	@Failure		400		{object}	ErrorResponse
//	@Failure		401		{object}	ErrorResponse
//	@Failure		403		{object}	ErrorResponse
//	@Failure		409		{object}	ErrorResponse
//	@Security		BearerAuth
//	@Router			/auth/invites [post]
func (h *Handler) CreateInvite(w http.ResponseWriter, r *http.Request) {
	user := UserFromContext(r.Context())
	if user == nil {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req CreateInviteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Email == "" {
		respondError(w, http.StatusBadRequest, "email is required")
		return
	}
	if req.Role == "" {
		respondError(w, http.StatusBadRequest, "role is required")
		return
	}
	if req.Role != RoleAdmin && req.Role != RoleEditor && req.Role != RoleViewer {
		respondError(w, http.StatusBadRequest, "invalid role")
		return
	}

	resp, err := h.service.CreateInvite(r.Context(), req, user)
	if err != nil {
		if errors.Is(err, ErrUserAlreadyExists) {
			respondError(w, http.StatusConflict, "user with this email already exists")
			return
		}
		if errors.Is(err, ErrPendingInviteExists) {
			respondError(w, http.StatusConflict, "a pending invite already exists for this email")
			return
		}
		log.Printf("create invite error: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to create invite")
		return
	}

	respondJSON(w, http.StatusCreated, resp)
}

// ListInvites lists pending invites.
//
//	@Summary		List invites
//	@Description	List pending invites for the organization
//	@Tags			auth
//	@Produce		json
//	@Success		200	{object}	ListInvitesResponse
//	@Failure		401	{object}	ErrorResponse
//	@Failure		403	{object}	ErrorResponse
//	@Security		BearerAuth
//	@Router			/auth/invites [get]
func (h *Handler) ListInvites(w http.ResponseWriter, r *http.Request) {
	user := UserFromContext(r.Context())
	if user == nil {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	invites, err := h.service.ListInvites(r.Context(), user.OrganizationID)
	if err != nil {
		log.Printf("list invites error: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to list invites")
		return
	}

	respondJSON(w, http.StatusOK, ListInvitesResponse{Invites: invites})
}

// CancelInvite cancels a pending invite.
//
//	@Summary		Cancel invite
//	@Description	Cancel a pending invite
//	@Tags			auth
//	@Produce		json
//	@Param			id	path		string	true	"Invite ID"
//	@Success		200	{object}	MessageResponse
//	@Failure		401	{object}	ErrorResponse
//	@Failure		403	{object}	ErrorResponse
//	@Failure		404	{object}	ErrorResponse
//	@Security		BearerAuth
//	@Router			/auth/invites/{id} [delete]
func (h *Handler) CancelInvite(w http.ResponseWriter, r *http.Request) {
	user := UserFromContext(r.Context())
	if user == nil {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	inviteID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid invite id")
		return
	}

	err = h.service.CancelInvite(r.Context(), inviteID, user.OrganizationID)
	if err != nil {
		if errors.Is(err, ErrInviteNotFound) {
			respondError(w, http.StatusNotFound, "invite not found")
			return
		}
		if errors.Is(err, ErrInviteAlreadyUsed) {
			respondError(w, http.StatusBadRequest, "invite has already been accepted")
			return
		}
		log.Printf("cancel invite error: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to cancel invite")
		return
	}

	respondJSON(w, http.StatusOK, MessageResponse{Message: "Invite cancelled"})
}

// ValidateInvite validates an invite token.
//
//	@Summary		Validate invite
//	@Description	Validate an invite token and get invite info
//	@Tags			auth
//	@Produce		json
//	@Param			token	query		string	true	"Invite token"
//	@Success		200		{object}	ValidateInviteResponse
//	@Failure		400		{object}	ErrorResponse
//	@Router			/auth/invites/validate [get]
func (h *Handler) ValidateInvite(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		respondError(w, http.StatusBadRequest, "token is required")
		return
	}

	resp, err := h.service.ValidateInvite(r.Context(), token)
	if err != nil {
		log.Printf("validate invite error: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to validate invite")
		return
	}

	respondJSON(w, http.StatusOK, resp)
}

// AcceptInvite accepts an invite and creates a user.
//
//	@Summary		Accept invite
//	@Description	Accept an invite and create a user account
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		AcceptInviteRequest	true	"Accept invite data"
//	@Success		201		{object}	MessageResponse
//	@Failure		400		{object}	ErrorResponse
//	@Failure		409		{object}	ErrorResponse
//	@Router			/auth/invites/accept [post]
func (h *Handler) AcceptInvite(w http.ResponseWriter, r *http.Request) {
	var req AcceptInviteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Token == "" {
		respondError(w, http.StatusBadRequest, "token is required")
		return
	}
	if req.Name == "" {
		respondError(w, http.StatusBadRequest, "name is required")
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

	_, err := h.service.AcceptInvite(r.Context(), req)
	if err != nil {
		if errors.Is(err, ErrInviteNotFound) {
			respondError(w, http.StatusBadRequest, "invalid invite token")
			return
		}
		if errors.Is(err, ErrInviteExpired) {
			respondError(w, http.StatusBadRequest, "invite has expired")
			return
		}
		if errors.Is(err, ErrInviteAlreadyUsed) {
			respondError(w, http.StatusBadRequest, "invite has already been used")
			return
		}
		if errors.Is(err, ErrEmailAlreadyExists) {
			respondError(w, http.StatusConflict, "email already exists")
			return
		}
		log.Printf("accept invite error: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to accept invite")
		return
	}

	respondJSON(w, http.StatusCreated, MessageResponse{Message: "Account created successfully. You can now log in."})
}

// ListUsers lists organization users.
//
//	@Summary		List users
//	@Description	List all users in the organization
//	@Tags			auth
//	@Produce		json
//	@Success		200	{object}	ListUsersResponse
//	@Failure		401	{object}	ErrorResponse
//	@Security		BearerAuth
//	@Router			/auth/users [get]
func (h *Handler) ListUsers(w http.ResponseWriter, r *http.Request) {
	user := UserFromContext(r.Context())
	if user == nil {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	users, err := h.service.ListUsers(r.Context(), user.OrganizationID)
	if err != nil {
		log.Printf("list users error: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to list users")
		return
	}

	respondJSON(w, http.StatusOK, ListUsersResponse{Users: users})
}

// UpdateUserRole updates a user's role.
//
//	@Summary		Update user role
//	@Description	Update a user's role in the organization
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string					true	"User ID"
//	@Param			request	body		UpdateUserRoleRequest	true	"Role data"
//	@Success		200		{object}	MessageResponse
//	@Failure		400		{object}	ErrorResponse
//	@Failure		401		{object}	ErrorResponse
//	@Failure		403		{object}	ErrorResponse
//	@Failure		404		{object}	ErrorResponse
//	@Security		BearerAuth
//	@Router			/auth/users/{id}/role [patch]
func (h *Handler) UpdateUserRole(w http.ResponseWriter, r *http.Request) {
	user := UserFromContext(r.Context())
	if user == nil {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	userID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	var req UpdateUserRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Role == "" {
		respondError(w, http.StatusBadRequest, "role is required")
		return
	}
	if req.Role != RoleAdmin && req.Role != RoleEditor && req.Role != RoleViewer {
		respondError(w, http.StatusBadRequest, "invalid role")
		return
	}

	err = h.service.UpdateUserRole(r.Context(), userID, req.Role, user)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			respondError(w, http.StatusNotFound, "user not found")
			return
		}
		if errors.Is(err, ErrLastAdmin) {
			respondError(w, http.StatusBadRequest, "cannot demote the last admin")
			return
		}
		log.Printf("update user role error: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to update user role")
		return
	}

	respondJSON(w, http.StatusOK, MessageResponse{Message: "User role updated"})
}

// RemoveUser removes a user from the organization.
//
//	@Summary		Remove user
//	@Description	Remove a user from the organization
//	@Tags			auth
//	@Produce		json
//	@Param			id	path		string	true	"User ID"
//	@Success		200	{object}	MessageResponse
//	@Failure		400		{object}	ErrorResponse
//	@Failure		401		{object}	ErrorResponse
//	@Failure		403		{object}	ErrorResponse
//	@Failure		404		{object}	ErrorResponse
//	@Security		BearerAuth
//	@Router			/auth/users/{id} [delete]
func (h *Handler) RemoveUser(w http.ResponseWriter, r *http.Request) {
	user := UserFromContext(r.Context())
	if user == nil {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	userID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	err = h.service.RemoveUser(r.Context(), userID, user)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			respondError(w, http.StatusNotFound, "user not found")
			return
		}
		if errors.Is(err, ErrCannotRemoveSelf) {
			respondError(w, http.StatusBadRequest, "cannot remove yourself")
			return
		}
		if errors.Is(err, ErrLastAdmin) {
			respondError(w, http.StatusBadRequest, "cannot remove the last admin")
			return
		}
		log.Printf("remove user error: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to remove user")
		return
	}

	respondJSON(w, http.StatusOK, MessageResponse{Message: "User removed"})
}

func (h *Handler) redirectWithOAuthResult(w http.ResponseWriter, r *http.Request, result *OAuthResult) {
	if result.PendingSetup != nil {
		// Redirect to complete-setup page
		redirectURL := h.service.GetAppURL() + "/auth/complete-setup" +
			"?token=" + url.QueryEscape(result.PendingSetup.Token) +
			"&email=" + url.QueryEscape(result.PendingSetup.Email) +
			"&name=" + url.QueryEscape(result.PendingSetup.ProviderName)
		http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
		return
	}

	// Existing user - redirect with auth
	h.redirectWithAuth(w, r, result.AuthResponse)
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
