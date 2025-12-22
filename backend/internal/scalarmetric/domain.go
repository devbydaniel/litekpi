package scalarmetric

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Error definitions
var (
	ErrScalarMetricNotFound  = errors.New("scalar metric not found")
	ErrUnauthorized          = errors.New("unauthorized")
	ErrLabelEmpty            = errors.New("label is required")
	ErrLabelTooLong          = errors.New("label exceeds maximum length of 255 characters")
	ErrMeasurementNameEmpty  = errors.New("measurement name is required")
	ErrInvalidTimeframe      = errors.New("invalid timeframe")
	ErrInvalidAggregation    = errors.New("invalid aggregation type")
	ErrInvalidComparisonType = errors.New("invalid comparison display type")
)

// Aggregation represents the aggregation type for a scalar metric.
type Aggregation string

const (
	AggregationSum     Aggregation = "sum"
	AggregationAverage Aggregation = "average"
)

// IsValid checks if the aggregation is valid.
func (a Aggregation) IsValid() bool {
	switch a {
	case AggregationSum, AggregationAverage:
		return true
	}
	return false
}

// ComparisonDisplayType represents how to display comparison values.
type ComparisonDisplayType string

const (
	ComparisonDisplayTypePercent  ComparisonDisplayType = "percent"
	ComparisonDisplayTypeAbsolute ComparisonDisplayType = "absolute"
)

// IsValid checks if the comparison display type is valid.
func (c ComparisonDisplayType) IsValid() bool {
	switch c {
	case ComparisonDisplayTypePercent, ComparisonDisplayTypeAbsolute:
		return true
	}
	return false
}

// Valid timeframes
var validTimeframes = map[string]bool{
	"last_7_days":  true,
	"last_30_days": true,
	"this_month":   true,
	"last_month":   true,
}

// IsValidTimeframe checks if the timeframe is valid.
func IsValidTimeframe(timeframe string) bool {
	return validTimeframes[timeframe]
}

// Filter represents a metadata filter for a scalar metric.
type Filter struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// ScalarMetric represents a scalar metric that belongs to a dashboard.
type ScalarMetric struct {
	ID          uuid.UUID `json:"id"`
	DashboardID uuid.UUID `json:"dashboardId"`
	Label       string    `json:"label"`
	Position    int       `json:"position"`

	// Query fields
	DataSourceID    uuid.UUID `json:"dataSourceId"`
	MeasurementName string    `json:"measurementName"`
	Timeframe       string    `json:"timeframe"` // last_7_days, last_30_days, this_month, last_month
	Filters         []Filter  `json:"filters"`

	// Calculation
	Aggregation Aggregation `json:"aggregation"`

	// Display fields
	ComparisonEnabled     bool                   `json:"comparisonEnabled"`
	ComparisonDisplayType *ComparisonDisplayType `json:"comparisonDisplayType,omitempty"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// ComputedScalarMetric represents a scalar metric with its calculated values.
type ComputedScalarMetric struct {
	ScalarMetric
	Value         float64  `json:"value"`
	PreviousValue *float64 `json:"previousValue,omitempty"`
	Change        *float64 `json:"change,omitempty"`
	ChangePercent *float64 `json:"changePercent,omitempty"`
}

// Request/Response types

// CreateScalarMetricRequest is the request body for creating a scalar metric.
type CreateScalarMetricRequest struct {
	DataSourceID          uuid.UUID              `json:"dataSourceId"`
	Label                 string                 `json:"label"`
	MeasurementName       string                 `json:"measurementName"`
	Timeframe             string                 `json:"timeframe"`
	Aggregation           Aggregation            `json:"aggregation"`
	Filters               []Filter               `json:"filters,omitempty"`
	ComparisonEnabled     bool                   `json:"comparisonEnabled"`
	ComparisonDisplayType *ComparisonDisplayType `json:"comparisonDisplayType,omitempty"`
}

// UpdateScalarMetricRequest is the request body for updating a scalar metric.
type UpdateScalarMetricRequest struct {
	Label                 string                 `json:"label"`
	Timeframe             string                 `json:"timeframe"`
	Aggregation           Aggregation            `json:"aggregation"`
	Filters               []Filter               `json:"filters,omitempty"`
	ComparisonEnabled     bool                   `json:"comparisonEnabled"`
	ComparisonDisplayType *ComparisonDisplayType `json:"comparisonDisplayType,omitempty"`
}

// ReorderScalarMetricsRequest is the request body for reordering scalar metrics.
type ReorderScalarMetricsRequest struct {
	ScalarMetricIDs []uuid.UUID `json:"scalarMetricIds"`
}

// ListScalarMetricsResponse is the response for listing scalar metrics.
type ListScalarMetricsResponse struct {
	ScalarMetrics []ScalarMetric `json:"scalarMetrics"`
}

// ComputeScalarMetricsResponse is the response for computing scalar metrics.
type ComputeScalarMetricsResponse struct {
	ScalarMetrics []ComputedScalarMetric `json:"scalarMetrics"`
}

// MessageResponse is a generic response with a message.
type MessageResponse struct {
	Message string `json:"message"`
}

// ErrorResponse represents an API error.
type ErrorResponse struct {
	Error string `json:"error"`
}
