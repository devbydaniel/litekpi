package email

import (
	"fmt"
	"net/smtp"
)

// Config holds email service configuration.
type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	From     string
}

// Service handles sending emails.
type Service struct {
	host     string
	port     int
	user     string
	password string
	from     string
	enabled  bool
}

// NewService creates a new email service.
func NewService(cfg Config) *Service {
	enabled := cfg.Host != "" && cfg.From != ""
	return &Service{
		host:     cfg.Host,
		port:     cfg.Port,
		user:     cfg.User,
		password: cfg.Password,
		from:     cfg.From,
		enabled:  enabled,
	}
}

// IsEnabled returns true if email service is configured.
func (s *Service) IsEnabled() bool {
	return s.enabled
}

// Send sends an email with the given recipient, subject, and body.
func (s *Service) Send(to, subject, body string) error {
	if !s.enabled {
		return nil // Silently skip if email not configured
	}

	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s",
		s.from, to, subject, body)

	addr := fmt.Sprintf("%s:%d", s.host, s.port)

	var auth smtp.Auth
	if s.user != "" && s.password != "" {
		auth = smtp.PlainAuth("", s.user, s.password, s.host)
	}

	return smtp.SendMail(addr, auth, s.from, []string{to}, []byte(msg))
}
