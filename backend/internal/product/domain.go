package product

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Product represents a product in the system.
type Product struct {
	ID             uuid.UUID `json:"id"`
	Name           string    `json:"name"`
	OrganizationID uuid.UUID `json:"organizationId"`
	APIKeyHash     string    `json:"-"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

// Error definitions
var (
	ErrProductNotFound  = errors.New("product not found")
	ErrUnauthorized     = errors.New("unauthorized")
	ErrProductNameEmpty = errors.New("product name is required")
)

// CreateProductRequest is the request body for creating a product.
type CreateProductRequest struct {
	Name string `json:"name"`
}

// CreateProductResponse is the response body for product creation.
type CreateProductResponse struct {
	Product Product `json:"product"`
	APIKey  string  `json:"apiKey"`
}

// RegenerateKeyResponse is the response body for API key regeneration.
type RegenerateKeyResponse struct {
	APIKey string `json:"apiKey"`
}

// ListProductsResponse is the response body for listing products.
type ListProductsResponse struct {
	Products []Product `json:"products"`
}

// MessageResponse is a generic response with a message.
type MessageResponse struct {
	Message string `json:"message"`
}

// ErrorResponse represents an API error.
type ErrorResponse struct {
	Error string `json:"error"`
}
