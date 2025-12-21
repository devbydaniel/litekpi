package report

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/devbydaniel/litekpi/internal/kpi"
)

// Error definitions
var (
	ErrReportNotFound  = errors.New("report not found")
	ErrUnauthorized    = errors.New("unauthorized")
	ErrReportNameEmpty = errors.New("report name is required")
	ErrReportNameTooLong = errors.New("report name exceeds maximum length of 255 characters")
)

// Report represents a report containing KPIs.
type Report struct {
	ID             uuid.UUID `json:"id"`
	Name           string    `json:"name"`
	OrganizationID uuid.UUID `json:"organizationId"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

// ReportWithKPIs is a report with its KPIs.
type ReportWithKPIs struct {
	Report Report    `json:"report"`
	KPIs   []kpi.KPI `json:"kpis"`
}

// ComputedReport is a report with computed KPI values.
type ComputedReport struct {
	Report Report            `json:"report"`
	KPIs   []kpi.ComputedKPI `json:"kpis"`
}

// Request/Response types

// CreateReportRequest is the request body for creating a report.
type CreateReportRequest struct {
	Name string `json:"name"`
}

// UpdateReportRequest is the request body for updating a report.
type UpdateReportRequest struct {
	Name string `json:"name"`
}

// ListReportsResponse is the response for listing reports.
type ListReportsResponse struct {
	Reports []Report `json:"reports"`
}

// MessageResponse is a generic response with a message.
type MessageResponse struct {
	Message string `json:"message"`
}

// ErrorResponse represents an API error.
type ErrorResponse struct {
	Error string `json:"error"`
}
