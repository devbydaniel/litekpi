package datasource

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// DataSource represents a data source in the system.
type DataSource struct {
	ID             uuid.UUID `json:"id"`
	Name           string    `json:"name"`
	OrganizationID uuid.UUID `json:"organizationId"`
	APIKeyHash     string    `json:"-"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

// Error definitions
var (
	ErrDataSourceNotFound  = errors.New("data source not found")
	ErrUnauthorized        = errors.New("unauthorized")
	ErrDataSourceNameEmpty = errors.New("data source name is required")
)

// CreateDataSourceRequest is the request body for creating a data source.
type CreateDataSourceRequest struct {
	Name string `json:"name"`
}

// CreateDataSourceResponse is the response body for data source creation.
type CreateDataSourceResponse struct {
	DataSource DataSource `json:"dataSource"`
	APIKey     string     `json:"apiKey"`
}

// RegenerateKeyResponse is the response body for API key regeneration.
type RegenerateKeyResponse struct {
	APIKey string `json:"apiKey"`
}

// ListDataSourcesResponse is the response body for listing data sources.
type ListDataSourcesResponse struct {
	DataSources []DataSource `json:"dataSources"`
}

// MessageResponse is a generic response with a message.
type MessageResponse struct {
	Message string `json:"message"`
}

// ErrorResponse represents an API error.
type ErrorResponse struct {
	Error string `json:"error"`
}
