package mcp

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// MCPAPIKey represents an MCP API key for organization-level access.
type MCPAPIKey struct {
	ID                   uuid.UUID   `json:"id"`
	OrganizationID       uuid.UUID   `json:"organizationId"`
	Name                 string      `json:"name"`
	APIKeyHash           string      `json:"-"`
	CreatedBy            uuid.UUID   `json:"createdBy"`
	LastUsedAt           *time.Time  `json:"lastUsedAt,omitempty"`
	CreatedAt            time.Time   `json:"createdAt"`
	AllowedDataSourceIDs []uuid.UUID `json:"allowedDataSourceIds"`
}

// Error definitions
var (
	ErrKeyNotFound           = errors.New("MCP API key not found")
	ErrUnauthorized          = errors.New("unauthorized")
	ErrKeyNameEmpty          = errors.New("API key name is required")
	ErrNoDataSourcesSelected = errors.New("at least one data source must be selected")
	ErrInvalidDataSource     = errors.New("invalid or unauthorized data source")
)

// CreateKeyRequest is the request body for creating an MCP API key.
type CreateKeyRequest struct {
	Name          string      `json:"name" validate:"required,max=255"`
	DataSourceIDs []uuid.UUID `json:"dataSourceIds" validate:"required,min=1"`
}

// UpdateKeyRequest is the request body for updating an MCP API key's data sources.
type UpdateKeyRequest struct {
	DataSourceIDs []uuid.UUID `json:"dataSourceIds" validate:"required,min=1"`
}

// CreateKeyResponse is the response body for API key creation.
type CreateKeyResponse struct {
	Key    MCPAPIKey `json:"key"`
	APIKey string    `json:"apiKey"` // Plain key, shown only once
}

// ListKeysResponse is the response body for listing MCP API keys.
type ListKeysResponse struct {
	Keys []MCPAPIKey `json:"keys"`
}

// MessageResponse is a generic response with a message.
type MessageResponse struct {
	Message string `json:"message"`
}

// ErrorResponse represents an API error.
type ErrorResponse struct {
	Error string `json:"error"`
}
