package config

import (
	"github.com/caarlos0/env/v11"
)

// Config holds all configuration for the application.
type Config struct {
	DatabaseURL string `env:"DATABASE_URL" envDefault:"postgres://trackable:secret@localhost:5432/trackable?sslmode=disable"`
	JWTSecret   string `env:"JWT_SECRET" envDefault:"your-secret-key-change-in-production"`
	AppURL      string `env:"APP_URL" envDefault:"http://localhost:5173"`
	APIURL      string `env:"API_URL" envDefault:"http://localhost:8080"`
	ServerPort  string `env:"SERVER_PORT" envDefault:"8080"`

	SMTP  SMTPConfig  `envPrefix:"SMTP_"`
	OAuth OAuthConfig `envPrefix:"OAUTH_"`
}

// SMTPConfig holds email configuration.
type SMTPConfig struct {
	Host     string `env:"HOST"`
	Port     int    `env:"PORT" envDefault:"587"`
	User     string `env:"USER"`
	Password string `env:"PASSWORD"`
	From     string `env:"FROM"`
}

// OAuthConfig holds OAuth provider credentials.
type OAuthConfig struct {
	GoogleClientID     string `env:"GOOGLE_CLIENT_ID"`
	GoogleClientSecret string `env:"GOOGLE_CLIENT_SECRET"`
	GithubClientID     string `env:"GITHUB_CLIENT_ID"`
	GithubClientSecret string `env:"GITHUB_CLIENT_SECRET"`
}

// Load reads configuration from environment variables.
func Load() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
