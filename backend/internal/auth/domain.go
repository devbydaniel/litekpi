package auth

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user in the system.
type User struct {
	ID            uuid.UUID  `json:"id"`
	Email         string     `json:"email"`
	PasswordHash  *string    `json:"-"`
	EmailVerified bool       `json:"emailVerified"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
}

// OAuthAccount represents a linked OAuth provider account.
type OAuthAccount struct {
	ID             uuid.UUID
	UserID         uuid.UUID
	Provider       string
	ProviderUserID string
	ProviderEmail  string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// EmailVerificationToken represents a token for email verification.
type EmailVerificationToken struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Token     string
	ExpiresAt time.Time
	CreatedAt time.Time
}

// PasswordResetToken represents a token for password reset.
type PasswordResetToken struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Token     string
	ExpiresAt time.Time
	Used      bool
	CreatedAt time.Time
}

// RegisterRequest is the request body for user registration.
type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginRequest is the request body for user login.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthResponse is the response body for successful authentication.
type AuthResponse struct {
	User  User   `json:"user"`
	Token string `json:"token"`
}

// ForgotPasswordRequest is the request body for initiating password reset.
type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

// ResetPasswordRequest is the request body for resetting password.
type ResetPasswordRequest struct {
	Token       string `json:"token"`
	NewPassword string `json:"newPassword"`
}

// VerifyEmailRequest is the request body for email verification.
type VerifyEmailRequest struct {
	Token string `json:"token"`
}

// ResendVerificationRequest is the request body for resending verification email.
type ResendVerificationRequest struct {
	Email string `json:"email"`
}

// MessageResponse is a generic response with a message.
type MessageResponse struct {
	Message string `json:"message"`
}

// ErrorResponse represents an API error.
type ErrorResponse struct {
	Error string `json:"error"`
}

// RegisterResponse is the response body for successful registration.
type RegisterResponse struct {
	Message string `json:"message"`
	User    User   `json:"user"`
}

// OAuthUserInfo represents user info from OAuth providers.
type OAuthUserInfo struct {
	ID    string
	Email string
	Name  string
}
