package ingest

import (
	"errors"
	"regexp"
	"time"

	"github.com/google/uuid"
)

// Validation constants
const (
	MaxMetricNameLength    = 128
	MaxBatchSize           = 100
	MaxMetadataKeys        = 20
	MaxMetadataKeyLength   = 64
	MaxMetadataValueLength = 256
)

// MetricNameRegex defines the valid pattern for metric names (snake_case).
var MetricNameRegex = regexp.MustCompile(`^[a-z][a-z0-9_]*$`)

// Error definitions
var (
	ErrInvalidAPIKey         = errors.New("invalid API key")
	ErrDuplicateMeasurement  = errors.New("duplicate measurement")
	ErrBatchTooLarge         = errors.New("batch exceeds maximum size")
	ErrEmptyBatch            = errors.New("batch must contain at least one measurement")
	ErrBatchDuplicates       = errors.New("batch contains duplicate measurements")
)

// Measurement represents a stored measurement data point.
type Measurement struct {
	ID        uuid.UUID          `json:"id"`
	ProductID uuid.UUID          `json:"productId"`
	Name      string             `json:"name"`
	Value     float64            `json:"value"`
	Timestamp time.Time          `json:"timestamp"`
	Metadata  map[string]string  `json:"metadata,omitempty"`
	CreatedAt time.Time          `json:"createdAt"`
}

// IngestRequest represents a single metric ingestion request.
type IngestRequest struct {
	Name      string            `json:"name"`
	Value     float64           `json:"value"`
	Timestamp string            `json:"timestamp,omitempty"`
	Metadata  map[string]string `json:"metadata,omitempty"`
}

// IngestResponse represents the response for a successful single metric ingestion.
type IngestResponse struct {
	ID        uuid.UUID         `json:"id"`
	Name      string            `json:"name"`
	Value     float64           `json:"value"`
	Timestamp time.Time         `json:"timestamp"`
	Metadata  map[string]string `json:"metadata,omitempty"`
}

// BatchIngestRequest represents a batch metric ingestion request.
type BatchIngestRequest struct {
	Metrics []IngestRequest `json:"metrics"`
}

// BatchIngestResponse represents the response for a successful batch ingestion.
type BatchIngestResponse struct {
	Count int `json:"count"`
}

// ValidationError represents an API validation error response.
type ValidationError struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// ErrorResponse represents a generic API error response.
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// MeasurementSummary represents a unique measurement name for a product.
type MeasurementSummary struct {
	Name         string   `json:"name"`
	MetadataKeys []string `json:"metadataKeys"`
}

// MetadataValues represents available values for a metadata key.
type MetadataValues struct {
	Key    string   `json:"key"`
	Values []string `json:"values"`
}

// AggregatedDataPoint represents a daily aggregated measurement.
type AggregatedDataPoint struct {
	Date  string  `json:"date"`  // YYYY-MM-DD format
	Sum   float64 `json:"sum"`   // Sum of values for the day
	Count int     `json:"count"` // Number of measurements
}

// ListMeasurementNamesResponse for listing unique measurement names.
type ListMeasurementNamesResponse struct {
	Measurements []MeasurementSummary `json:"measurements"`
}

// GetMetadataValuesResponse for metadata filter options.
type GetMetadataValuesResponse struct {
	Metadata []MetadataValues `json:"metadata"`
}

// GetMeasurementDataResponse for chart data.
type GetMeasurementDataResponse struct {
	Name       string                `json:"name"`
	DataPoints []AggregatedDataPoint `json:"dataPoints"`
}

// SplitSeries represents aggregated data for a single metadata value.
type SplitSeries struct {
	Key        string                `json:"key"`
	DataPoints []AggregatedDataPoint `json:"dataPoints"`
}

// GetMeasurementDataSplitResponse for chart data split by a metadata key.
type GetMeasurementDataSplitResponse struct {
	Name    string        `json:"name"`
	SplitBy string        `json:"splitBy"`
	Series  []SplitSeries `json:"series"`
}
