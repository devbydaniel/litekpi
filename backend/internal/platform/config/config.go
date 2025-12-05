package config

import (
	"os"
	"strconv"
)

// Config holds all configuration for the application.
type Config struct {
	DatabaseURL string
	JWTSecret   string
	AppURL      string
	APIURL      string
	ServerPort  string

	// SMTP settings
	SMTP SMTPConfig

	// OAuth settings
	OAuth OAuthConfig
}

// SMTPConfig holds email configuration.
type SMTPConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	From     string
}

// OAuthConfig holds OAuth provider credentials.
type OAuthConfig struct {
	GoogleClientID     string
	GoogleClientSecret string
	GithubClientID     string
	GithubClientSecret string
}

// Load reads configuration from environment variables.
func Load() *Config {
	return &Config{
		DatabaseURL: getEnv("DATABASE_URL", "postgres://trackable:secret@localhost:5432/trackable?sslmode=disable"),
		JWTSecret:   getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		AppURL:      getEnv("APP_URL", "http://localhost:5173"),
		APIURL:      getEnv("API_URL", "http://localhost:8080"),
		ServerPort:  getEnv("SERVER_PORT", "8080"),

		SMTP: SMTPConfig{
			Host:     getEnv("SMTP_HOST", ""),
			Port:     getEnvInt("SMTP_PORT", 587),
			User:     getEnv("SMTP_USER", ""),
			Password: getEnv("SMTP_PASSWORD", ""),
			From:     getEnv("SMTP_FROM", ""),
		},

		OAuth: OAuthConfig{
			GoogleClientID:     getEnv("OAUTH_GOOGLE_CLIENT_ID", ""),
			GoogleClientSecret: getEnv("OAUTH_GOOGLE_CLIENT_SECRET", ""),
			GithubClientID:     getEnv("OAUTH_GITHUB_CLIENT_ID", ""),
			GithubClientSecret: getEnv("OAUTH_GITHUB_CLIENT_SECRET", ""),
		},
	}
}

// getEnv returns the value of an environment variable or a default value.
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// getEnvInt returns the integer value of an environment variable or a default value.
func getEnvInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
