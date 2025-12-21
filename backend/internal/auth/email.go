package auth

import (
	"fmt"
	"net/smtp"
	"strings"
)

// EmailConfig holds email service configuration.
type EmailConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	From     string
	AppURL   string
}

// EmailService handles sending emails.
type EmailService struct {
	host     string
	port     int
	user     string
	password string
	from     string
	appURL   string
	enabled  bool
}

// NewEmailService creates a new email service.
func NewEmailService(cfg EmailConfig) *EmailService {
	enabled := cfg.Host != "" && cfg.From != ""
	return &EmailService{
		host:     cfg.Host,
		port:     cfg.Port,
		user:     cfg.User,
		password: cfg.Password,
		from:     cfg.From,
		appURL:   strings.TrimSuffix(cfg.AppURL, "/"),
		enabled:  enabled,
	}
}

// IsEnabled returns true if email service is configured.
func (e *EmailService) IsEnabled() bool {
	return e.enabled
}

// SendVerificationEmail sends an email verification link.
func (e *EmailService) SendVerificationEmail(to, token string) error {
	if !e.enabled {
		return nil // Silently skip if email not configured
	}

	subject := "Verify your LiteKPI account"
	verifyURL := fmt.Sprintf("%s/verify-email?token=%s", e.appURL, token)

	body := fmt.Sprintf(`Hi,

Thanks for signing up for LiteKPI! Please verify your email address by clicking the link below:

%s

This link will expire in 24 hours.

If you didn't create a LiteKPI account, you can safely ignore this email.

Thanks,
The LiteKPI Team`, verifyURL)

	return e.sendEmail(to, subject, body)
}

// SendPasswordResetEmail sends a password reset link.
func (e *EmailService) SendPasswordResetEmail(to, token string) error {
	if !e.enabled {
		return nil // Silently skip if email not configured
	}

	subject := "Reset your LiteKPI password"
	resetURL := fmt.Sprintf("%s/new-password?token=%s", e.appURL, token)

	body := fmt.Sprintf(`Hi,

We received a request to reset your password. Click the link below to create a new password:

%s

This link will expire in 1 hour.

If you didn't request a password reset, you can safely ignore this email.

Thanks,
The LiteKPI Team`, resetURL)

	return e.sendEmail(to, subject, body)
}

// SendInviteEmail sends an invitation email.
func (e *EmailService) SendInviteEmail(to, token, inviterName, orgName string) error {
	if !e.enabled {
		return nil // Silently skip if email not configured
	}

	subject := fmt.Sprintf("You've been invited to join %s on LiteKPI", orgName)
	inviteURL := fmt.Sprintf("%s/accept-invite?token=%s", e.appURL, token)

	body := fmt.Sprintf(`Hi,

%s has invited you to join %s on LiteKPI.

Click the link below to accept the invitation and create your account:

%s

This invitation will expire in 7 days.

Thanks,
The LiteKPI Team`, inviterName, orgName, inviteURL)

	return e.sendEmail(to, subject, body)
}

func (e *EmailService) sendEmail(to, subject, body string) error {
	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s",
		e.from, to, subject, body)

	addr := fmt.Sprintf("%s:%d", e.host, e.port)

	var auth smtp.Auth
	if e.user != "" && e.password != "" {
		auth = smtp.PlainAuth("", e.user, e.password, e.host)
	}

	return smtp.SendMail(addr, auth, e.from, []string{to}, []byte(msg))
}
