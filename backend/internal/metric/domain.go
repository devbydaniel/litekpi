package metric

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Error definitions
var (
	ErrMetricNotFound         = errors.New("metric not found")
	ErrUnauthorized           = errors.New("unauthorized")
	ErrLabelEmpty             = errors.New("label is required")
	ErrLabelTooLong           = errors.New("label exceeds maximum length of 255 characters")
	ErrMeasurementNameEmpty   = errors.New("measurement name is required")
	ErrInvalidTimeframe       = errors.New("invalid timeframe")
	ErrInvalidAggregation     = errors.New("invalid aggregation type")
	ErrInvalidGranularity     = errors.New("invalid granularity")
	ErrInvalidDisplayMode     = errors.New("invalid display mode")
	ErrInvalidChartType       = errors.New("invalid chart type")
	ErrInvalidComparisonType  = errors.New("invalid comparison display type")
	ErrAggregationKeyRequired = errors.New("aggregation_key is required for count_unique aggregation")
	ErrChartTypeRequired      = errors.New("chart_type is required for time_series display mode")
)

// DisplayMode represents how the metric is displayed.
type DisplayMode string

const (
	DisplayModeScalar     DisplayMode = "scalar"
	DisplayModeTimeSeries DisplayMode = "time_series"
)

// IsValid checks if the display mode is valid.
func (d DisplayMode) IsValid() bool {
	switch d {
	case DisplayModeScalar, DisplayModeTimeSeries:
		return true
	}
	return false
}

// Aggregation represents the aggregation type for a metric.
type Aggregation string

const (
	AggregationSum         Aggregation = "sum"
	AggregationAverage     Aggregation = "average"
	AggregationCount       Aggregation = "count"
	AggregationCountUnique Aggregation = "count_unique"
)

// IsValid checks if the aggregation is valid.
func (a Aggregation) IsValid() bool {
	switch a {
	case AggregationSum, AggregationAverage, AggregationCount, AggregationCountUnique:
		return true
	}
	return false
}

// RequiresAggregationKey returns true if the aggregation type needs an aggregation_key.
func (a Aggregation) RequiresAggregationKey() bool {
	return a == AggregationCountUnique
}

// Granularity represents the time granularity for aggregation.
type Granularity string

const (
	GranularityDaily   Granularity = "daily"
	GranularityWeekly  Granularity = "weekly"
	GranularityMonthly Granularity = "monthly"
)

// IsValid checks if the granularity is valid.
func (g Granularity) IsValid() bool {
	switch g {
	case GranularityDaily, GranularityWeekly, GranularityMonthly:
		return true
	}
	return false
}

// ChartType represents the type of chart for time series display.
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
	"custom":       true,
}

// IsValidTimeframe checks if the timeframe is valid.
func IsValidTimeframe(timeframe string) bool {
	return validTimeframes[timeframe]
}

// Filter represents a metadata filter for a metric.
type Filter struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// Metric represents a unified metric on a dashboard.
type Metric struct {
	ID          uuid.UUID `json:"id"`
	DashboardID uuid.UUID `json:"dashboardId"`
	Label       string    `json:"label"`
	Position    int       `json:"position"`

	// Query fields
	DataSourceID    uuid.UUID  `json:"dataSourceId"`
	MeasurementName string     `json:"measurementName"`
	Timeframe       string     `json:"timeframe"` // last_7_days, last_30_days, this_month, last_month, custom
	DateFrom        *time.Time `json:"dateFrom,omitempty"`
	DateTo          *time.Time `json:"dateTo,omitempty"`
	Filters         []Filter   `json:"filters"`

	// Aggregation
	Aggregation    Aggregation `json:"aggregation"`
	AggregationKey *string     `json:"aggregationKey,omitempty"` // Required for count_unique
	Granularity    Granularity `json:"granularity"`

	// Display mode
	DisplayMode DisplayMode `json:"displayMode"`

	// Scalar display options
	ComparisonEnabled     bool                   `json:"comparisonEnabled"`
	ComparisonDisplayType *ComparisonDisplayType `json:"comparisonDisplayType,omitempty"`

	// Time series display options
	ChartType *ChartType `json:"chartType,omitempty"`
	SplitBy   *string    `json:"splitBy,omitempty"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// ComputedMetric represents a metric with its calculated values.
type ComputedMetric struct {
	Metric
	// For scalar display
	Value         *float64 `json:"value,omitempty"`
	PreviousValue *float64 `json:"previousValue,omitempty"`
	Change        *float64 `json:"change,omitempty"`
	ChangePercent *float64 `json:"changePercent,omitempty"`

	// For time series display
	DataPoints []DataPoint   `json:"dataPoints,omitempty"`
	Series     []SplitSeries `json:"series,omitempty"` // When splitBy is used
}

// DataPoint represents a single aggregated data point.
type DataPoint struct {
	Date  string  `json:"date"`
	Value float64 `json:"value"`
}

// SplitSeries represents aggregated data for a single metadata value.
type SplitSeries struct {
	Key        string      `json:"key"`
	DataPoints []DataPoint `json:"dataPoints"`
}

// AggregatedDataPoint represents raw aggregated data from the database.
type AggregatedDataPoint struct {
	Date  string
	Sum   float64
	Count int
}

// Request/Response types

// CreateMetricRequest is the request body for creating a metric.
type CreateMetricRequest struct {
	DataSourceID    uuid.UUID   `json:"dataSourceId"`
	Label           string      `json:"label"`
	MeasurementName string      `json:"measurementName"`
	Timeframe       string      `json:"timeframe"`
	DateFrom        *time.Time  `json:"dateFrom,omitempty"`
	DateTo          *time.Time  `json:"dateTo,omitempty"`
	Filters         []Filter    `json:"filters,omitempty"`
	Aggregation     Aggregation `json:"aggregation"`
	AggregationKey  *string     `json:"aggregationKey,omitempty"`
	Granularity     Granularity `json:"granularity"`
	DisplayMode     DisplayMode `json:"displayMode"`

	// Scalar options
	ComparisonEnabled     bool                   `json:"comparisonEnabled"`
	ComparisonDisplayType *ComparisonDisplayType `json:"comparisonDisplayType,omitempty"`

	// Time series options
	ChartType *ChartType `json:"chartType,omitempty"`
	SplitBy   *string    `json:"splitBy,omitempty"`
}

// UpdateMetricRequest is the request body for updating a metric.
type UpdateMetricRequest struct {
	Label           string      `json:"label"`
	Timeframe       string      `json:"timeframe"`
	DateFrom        *time.Time  `json:"dateFrom,omitempty"`
	DateTo          *time.Time  `json:"dateTo,omitempty"`
	Filters         []Filter    `json:"filters,omitempty"`
	Aggregation     Aggregation `json:"aggregation"`
	AggregationKey  *string     `json:"aggregationKey,omitempty"`
	Granularity     Granularity `json:"granularity"`
	DisplayMode     DisplayMode `json:"displayMode"`

	// Scalar options
	ComparisonEnabled     bool                   `json:"comparisonEnabled"`
	ComparisonDisplayType *ComparisonDisplayType `json:"comparisonDisplayType,omitempty"`

	// Time series options
	ChartType *ChartType `json:"chartType,omitempty"`
	SplitBy   *string    `json:"splitBy,omitempty"`
}

// ReorderMetricsRequest is the request body for reordering metrics.
type ReorderMetricsRequest struct {
	MetricIDs []uuid.UUID `json:"metricIds"`
}

// ListMetricsResponse is the response for listing metrics.
type ListMetricsResponse struct {
	Metrics []Metric `json:"metrics"`
}

// ComputeMetricsResponse is the response for computing metrics.
type ComputeMetricsResponse struct {
	Metrics []ComputedMetric `json:"metrics"`
}

// MessageResponse is a generic response with a message.
type MessageResponse struct {
	Message string `json:"message"`
}

// ErrorResponse represents an API error.
type ErrorResponse struct {
	Error string `json:"error"`
}
