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
	ErrInvalidAPIKey    = errors.New("invalid API key")
	ErrDuplicateMetric  = errors.New("duplicate metric")
	ErrBatchTooLarge    = errors.New("batch exceeds maximum size")
	ErrEmptyBatch       = errors.New("batch must contain at least one metric")
	ErrBatchDuplicates  = errors.New("batch contains duplicate metrics")
)

// Metric represents a stored metric data point.
type Metric struct {
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
