package auth

import (
	"fmt"
	"strings"

	"github.com/devbydaniel/litekpi/internal/platform/email"
)

// AuthEmailer handles sending auth-related emails.
type AuthEmailer struct {
	svc    *email.Service
	appURL string
}

// NewAuthEmailer creates a new auth emailer.
func NewAuthEmailer(svc *email.Service, appURL string) *AuthEmailer {
	return &AuthEmailer{
		svc:    svc,
		appURL: strings.TrimSuffix(appURL, "/"),
	}
}

// IsEnabled returns true if email service is configured.
func (e *AuthEmailer) IsEnabled() bool {
	return e.svc.IsEnabled()
}

// SendVerificationEmail sends an email verification link.
func (e *AuthEmailer) SendVerificationEmail(to, token string) error {
	subject := "Verify your LiteKPI account"
	verifyURL := fmt.Sprintf("%s/verify-email?token=%s", e.appURL, token)

	body := fmt.Sprintf(`Hi,

Thanks for signing up for LiteKPI! Please verify your email address by clicking the link below:

%s

This link will expire in 24 hours.

If you didn't create a LiteKPI account, you can safely ignore this email.

Thanks,
The LiteKPI Team`, verifyURL)

	return e.svc.Send(to, subject, body)
}

// SendPasswordResetEmail sends a password reset link.
func (e *AuthEmailer) SendPasswordResetEmail(to, token string) error {
	subject := "Reset your LiteKPI password"
	resetURL := fmt.Sprintf("%s/new-password?token=%s", e.appURL, token)

	body := fmt.Sprintf(`Hi,

We received a request to reset your password. Click the link below to create a new password:

%s

This link will expire in 1 hour.

If you didn't request a password reset, you can safely ignore this email.

Thanks,
The LiteKPI Team`, resetURL)

	return e.svc.Send(to, subject, body)
}

// SendInviteEmail sends an invitation email.
func (e *AuthEmailer) SendInviteEmail(to, token, inviterName, orgName string) error {
	subject := fmt.Sprintf("You've been invited to join %s on LiteKPI", orgName)
	inviteURL := fmt.Sprintf("%s/accept-invite?token=%s", e.appURL, token)

	body := fmt.Sprintf(`Hi,

%s has invited you to join %s on LiteKPI.

Click the link below to accept the invitation and create your account:

%s

This invitation will expire in 7 days.

Thanks,
The LiteKPI Team`, inviterName, orgName, inviteURL)

	return e.svc.Send(to, subject, body)
}
