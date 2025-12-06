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

	subject := "Verify your Trackable account"
	verifyURL := fmt.Sprintf("%s/verify-email?token=%s", e.appURL, token)

	body := fmt.Sprintf(`Hi,

Thanks for signing up for Trackable! Please verify your email address by clicking the link below:

%s

This link will expire in 24 hours.

If you didn't create a Trackable account, you can safely ignore this email.

Thanks,
The Trackable Team`, verifyURL)

	return e.sendEmail(to, subject, body)
}

// SendPasswordResetEmail sends a password reset link.
func (e *EmailService) SendPasswordResetEmail(to, token string) error {
	if !e.enabled {
		return nil // Silently skip if email not configured
	}

	subject := "Reset your Trackable password"
	resetURL := fmt.Sprintf("%s/new-password?token=%s", e.appURL, token)

	body := fmt.Sprintf(`Hi,

We received a request to reset your password. Click the link below to create a new password:

%s

This link will expire in 1 hour.

If you didn't request a password reset, you can safely ignore this email.

Thanks,
The Trackable Team`, resetURL)

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
