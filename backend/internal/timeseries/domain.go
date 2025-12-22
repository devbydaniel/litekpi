package timeseries

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Error definitions
var (
	ErrTimeSeriesNotFound   = errors.New("time series not found")
	ErrUnauthorized         = errors.New("unauthorized")
	ErrMeasurementNameEmpty = errors.New("measurement name is required")
	ErrInvalidChartType     = errors.New("invalid chart type")
	ErrInvalidDateRange     = errors.New("invalid date range")
	ErrTitleTooLong         = errors.New("title exceeds maximum length of 128 characters")
)

// ChartType represents the type of chart for a time series.
type ChartType string

const (
	ChartTypeArea ChartType = "area"
	ChartTypeBar  ChartType = "bar"
	ChartTypeLine ChartType = "line"
)

// IsValid checks if the chart type is valid.
func (c ChartType) IsValid() bool {
	switch c {
	case ChartTypeArea, ChartTypeBar, ChartTypeLine:
		return true
	}
	return false
}

// Valid date ranges
var validDateRanges = map[string]bool{
	"last_7_days":  true,
	"last_30_days": true,
	"custom":       true,
}

// IsValidDateRange checks if the date range is valid.
func IsValidDateRange(dateRange string) bool {
	return validDateRanges[dateRange]
}

// Filter represents a metadata filter for a time series.
type Filter struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// TimeSeries represents a time series chart on a dashboard.
type TimeSeries struct {
	ID          uuid.UUID  `json:"id"`
	DashboardID uuid.UUID  `json:"dashboardId"`
	Title       *string    `json:"title,omitempty"`
	Position    int        `json:"position"`

	// Query fields
	DataSourceID    uuid.UUID  `json:"dataSourceId"`
	MeasurementName string     `json:"measurementName"`
	DateRange       string     `json:"dateRange"` // last_7_days, last_30_days, custom
	DateFrom        *time.Time `json:"dateFrom,omitempty"`
	DateTo          *time.Time `json:"dateTo,omitempty"`
	SplitBy         *string    `json:"splitBy,omitempty"`
	Filters         []Filter   `json:"filters"`

	// Display fields
	ChartType ChartType `json:"chartType"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Request/Response types

// CreateTimeSeriesRequest is the request body for creating a time series.
type CreateTimeSeriesRequest struct {
	DataSourceID    uuid.UUID  `json:"dataSourceId"`
	MeasurementName string     `json:"measurementName"`
	Title           *string    `json:"title,omitempty"`
	ChartType       ChartType  `json:"chartType"`
	DateRange       string     `json:"dateRange"`
	DateFrom        *time.Time `json:"dateFrom,omitempty"`
	DateTo          *time.Time `json:"dateTo,omitempty"`
	SplitBy         *string    `json:"splitBy,omitempty"`
	Filters         []Filter   `json:"filters,omitempty"`
}

// UpdateTimeSeriesRequest is the request body for updating a time series.
type UpdateTimeSeriesRequest struct {
	Title     *string    `json:"title,omitempty"`
	ChartType ChartType  `json:"chartType"`
	DateRange string     `json:"dateRange"`
	DateFrom  *time.Time `json:"dateFrom,omitempty"`
	DateTo    *time.Time `json:"dateTo,omitempty"`
	SplitBy   *string    `json:"splitBy,omitempty"`
	Filters   []Filter   `json:"filters,omitempty"`
}

// ReorderTimeSeriesRequest is the request body for reordering time series.
type ReorderTimeSeriesRequest struct {
	TimeSeriesIDs []uuid.UUID `json:"timeSeriesIds"`
}

// ListTimeSeriesResponse is the response for listing time series.
type ListTimeSeriesResponse struct {
	TimeSeries []TimeSeries `json:"timeSeries"`
}

// MessageResponse is a generic response with a message.
type MessageResponse struct {
	Message string `json:"message"`
}

// ErrorResponse represents an API error.
type ErrorResponse struct {
	Error string `json:"error"`
}
