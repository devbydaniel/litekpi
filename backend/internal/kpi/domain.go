package kpi

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Error definitions
var (
	ErrKPINotFound           = errors.New("kpi not found")
	ErrUnauthorized          = errors.New("unauthorized")
	ErrLabelEmpty            = errors.New("label is required")
	ErrLabelTooLong          = errors.New("label exceeds maximum length of 255 characters")
	ErrMeasurementNameEmpty  = errors.New("measurement name is required")
	ErrInvalidTimeframe      = errors.New("invalid timeframe")
	ErrInvalidAggregation    = errors.New("invalid aggregation type")
	ErrInvalidComparisonType = errors.New("invalid comparison display type")
)

// KPI represents a KPI card that can belong to a dashboard or report.
type KPI struct {
	ID                    uuid.UUID  `json:"id"`
	DashboardID           *uuid.UUID `json:"dashboardId,omitempty"`
	ReportID              *uuid.UUID `json:"reportId,omitempty"`
	DataSourceID          uuid.UUID  `json:"dataSourceId"`
	Label                 string     `json:"label"`
	MeasurementName       string     `json:"measurementName"`
	Timeframe             string     `json:"timeframe"` // last_7_days, last_30_days, this_month, last_month
	Aggregation           string     `json:"aggregation"` // sum, average
	Filters               []Filter   `json:"filters"`
	ComparisonEnabled     bool       `json:"comparisonEnabled"`
	ComparisonDisplayType *string    `json:"comparisonDisplayType,omitempty"` // percent, absolute
	Position              int        `json:"position"`
	CreatedAt             time.Time  `json:"createdAt"`
	UpdatedAt             time.Time  `json:"updatedAt"`
}

// Filter represents a metadata filter for a KPI.
type Filter struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// ComputedKPI represents a KPI with its calculated values.
type ComputedKPI struct {
	KPI
	Value         float64  `json:"value"`
	PreviousValue *float64 `json:"previousValue,omitempty"`
	Change        *float64 `json:"change,omitempty"`
	ChangePercent *float64 `json:"changePercent,omitempty"`
}

// Valid timeframes
var validTimeframes = map[string]bool{
	"last_7_days":  true,
	"last_30_days": true,
	"this_month":   true,
	"last_month":   true,
}

// Valid aggregation types
var validAggregations = map[string]bool{
	"sum":     true,
	"average": true,
}

// Valid comparison display types
var validComparisonDisplayTypes = map[string]bool{
	"percent":  true,
	"absolute": true,
}

// IsValidTimeframe checks if the timeframe is valid.
func IsValidTimeframe(timeframe string) bool {
	return validTimeframes[timeframe]
}

// IsValidAggregation checks if the aggregation type is valid.
func IsValidAggregation(aggregation string) bool {
	return validAggregations[aggregation]
}

// IsValidComparisonDisplayType checks if the comparison display type is valid.
func IsValidComparisonDisplayType(displayType string) bool {
	return validComparisonDisplayTypes[displayType]
}

// Request/Response types

// CreateKPIRequest is the request body for creating a KPI.
type CreateKPIRequest struct {
	DataSourceID          uuid.UUID `json:"dataSourceId"`
	Label                 string    `json:"label"`
	MeasurementName       string    `json:"measurementName"`
	Timeframe             string    `json:"timeframe"`
	Aggregation           string    `json:"aggregation"`
	Filters               []Filter  `json:"filters,omitempty"`
	ComparisonEnabled     bool      `json:"comparisonEnabled"`
	ComparisonDisplayType *string   `json:"comparisonDisplayType,omitempty"`
}

// UpdateKPIRequest is the request body for updating a KPI.
type UpdateKPIRequest struct {
	Label                 string   `json:"label"`
	Timeframe             string   `json:"timeframe"`
	Aggregation           string   `json:"aggregation"`
	Filters               []Filter `json:"filters,omitempty"`
	ComparisonEnabled     bool     `json:"comparisonEnabled"`
	ComparisonDisplayType *string  `json:"comparisonDisplayType,omitempty"`
}

// ReorderKPIsRequest is the request body for reordering KPIs.
type ReorderKPIsRequest struct {
	KPIIDs []uuid.UUID `json:"kpiIds"`
}

// ListKPIsResponse is the response for listing KPIs.
type ListKPIsResponse struct {
	KPIs []KPI `json:"kpis"`
}

// ComputeKPIsResponse is the response for computing KPIs.
type ComputeKPIsResponse struct {
	KPIs []ComputedKPI `json:"kpis"`
}

// MessageResponse is a generic response with a message.
type MessageResponse struct {
	Message string `json:"message"`
}

// ErrorResponse represents an API error.
type ErrorResponse struct {
	Error string `json:"error"`
}
