package dashboard

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/devbydaniel/litekpi/internal/scalarmetric"
	"github.com/devbydaniel/litekpi/internal/timeseries"
)

// Error definitions
var (
	ErrDashboardNotFound   = errors.New("dashboard not found")
	ErrUnauthorized        = errors.New("unauthorized")
	ErrDashboardNameEmpty  = errors.New("dashboard name is required")
	ErrCannotDeleteDefault = errors.New("cannot delete default dashboard")
)

// Dashboard represents a dashboard in the system.
type Dashboard struct {
	ID             uuid.UUID `json:"id"`
	Name           string    `json:"name"`
	OrganizationID uuid.UUID `json:"organizationId"`
	IsDefault      bool      `json:"isDefault"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

// Request/Response types

// CreateDashboardRequest is the request body for creating a dashboard.
type CreateDashboardRequest struct {
	Name string `json:"name"`
}

// UpdateDashboardRequest is the request body for updating a dashboard.
type UpdateDashboardRequest struct {
	Name string `json:"name"`
}

// DashboardWithData is a dashboard with its time series and scalar metrics.
type DashboardWithData struct {
	Dashboard     Dashboard                `json:"dashboard"`
	TimeSeries    []timeseries.TimeSeries  `json:"timeSeries"`
	ScalarMetrics []scalarmetric.ScalarMetric `json:"scalarMetrics"`
}

// ListDashboardsResponse is the response for listing dashboards.
type ListDashboardsResponse struct {
	Dashboards []Dashboard `json:"dashboards"`
}

// MessageResponse is a generic response with a message.
type MessageResponse struct {
	Message string `json:"message"`
}

// ErrorResponse represents an API error.
type ErrorResponse struct {
	Error string `json:"error"`
}
