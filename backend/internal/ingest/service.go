package ingest

import (
	"context"
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/google/uuid"
)

// Service handles measurement ingestion business logic.
type Service struct {
	repo *Repository
}

// NewService creates a new ingest service.
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// IngestSingle validates and ingests a single measurement.
func (s *Service) IngestSingle(ctx context.Context, dataSourceID uuid.UUID, req IngestRequest) (*IngestResponse, error) {
	// Validate metric name
	if err := validateMetricName(req.Name); err != nil {
		return nil, err
	}

	// Validate value
	if err := validateValue(req.Value); err != nil {
		return nil, err
	}

	// Parse or default timestamp
	timestamp, err := parseTimestamp(req.Timestamp)
	if err != nil {
		return nil, err
	}

	// Validate metadata
	if err := validateMetadata(req.Metadata); err != nil {
		return nil, err
	}

	// Create measurement
	measurement, err := s.repo.CreateMeasurement(ctx, dataSourceID, req.Name, req.Value, timestamp, req.Metadata)
	if err != nil {
		return nil, err
	}

	return &IngestResponse{
		ID:        measurement.ID,
		Name:      measurement.Name,
		Value:     measurement.Value,
		Timestamp: measurement.Timestamp,
		Metadata:  measurement.Metadata,
	}, nil
}

// IngestBatch validates and ingests multiple measurements atomically.
func (s *Service) IngestBatch(ctx context.Context, dataSourceID uuid.UUID, req BatchIngestRequest) (*BatchIngestResponse, error) {
	// Validate batch size
	if len(req.Metrics) == 0 {
		return nil, &validationError{
			errorType: "validation_failed",
			message:   "Batch must contain at least one measurement",
		}
	}
	if len(req.Metrics) > MaxBatchSize {
		return nil, &validationError{
			errorType: "validation_failed",
			message:   fmt.Sprintf("Batch exceeds maximum size of %d measurements", MaxBatchSize),
		}
	}

	// Parse timestamps and check for internal duplicates
	timestamps := make([]time.Time, len(req.Metrics))
	seen := make(map[string]int) // key: "name|timestamp" -> index

	for i, m := range req.Metrics {
		// Validate metric name
		if err := validateMetricName(m.Name); err != nil {
			return nil, &validationError{
				errorType: "validation_failed",
				message:   fmt.Sprintf("Measurement at index %d: %s", i, err.Error()),
			}
		}

		// Validate value
		if err := validateValue(m.Value); err != nil {
			return nil, &validationError{
				errorType: "validation_failed",
				message:   fmt.Sprintf("Measurement at index %d: %s", i, err.Error()),
			}
		}

		// Parse timestamp
		ts, err := parseTimestamp(m.Timestamp)
		if err != nil {
			return nil, &validationError{
				errorType: "validation_failed",
				message:   fmt.Sprintf("Measurement at index %d: %s", i, err.Error()),
			}
		}
		timestamps[i] = ts

		// Validate metadata
		if err := validateMetadata(m.Metadata); err != nil {
			return nil, &validationError{
				errorType: "validation_failed",
				message:   fmt.Sprintf("Measurement at index %d: %s", i, err.Error()),
			}
		}

		// Check for internal duplicates
		key := fmt.Sprintf("%s|%s", m.Name, ts.Format(time.RFC3339Nano))
		if prevIdx, exists := seen[key]; exists {
			return nil, &validationError{
				errorType: "validation_failed",
				message:   fmt.Sprintf("Batch contains duplicate measurements (same name and timestamp) at indices %d and %d", prevIdx, i),
			}
		}
		seen[key] = i
	}

	// Insert all measurements
	count, err := s.repo.CreateMeasurementsBatch(ctx, dataSourceID, req.Metrics, timestamps)
	if err != nil {
		return nil, err
	}

	return &BatchIngestResponse{
		Count: count,
	}, nil
}

// validationError is a custom error type for validation failures.
type validationError struct {
	errorType string
	message   string
}

func (e *validationError) Error() string {
	return e.message
}

// IsValidationError checks if an error is a validation error.
func IsValidationError(err error) (*validationError, bool) {
	ve, ok := err.(*validationError)
	return ve, ok
}

// validateMetricName validates the metric name format.
func validateMetricName(name string) error {
	if name == "" {
		return &validationError{
			errorType: "validation_failed",
			message:   "metric name is required",
		}
	}
	if len(name) > MaxMetricNameLength {
		return &validationError{
			errorType: "validation_failed",
			message:   fmt.Sprintf("Invalid metric name '%s': exceeds maximum length of %d characters", name, MaxMetricNameLength),
		}
	}
	if !MetricNameRegex.MatchString(name) {
		return &validationError{
			errorType: "validation_failed",
			message:   fmt.Sprintf("Invalid metric name '%s': must be snake_case (lowercase alphanumeric and underscores, starting with letter)", name),
		}
	}
	return nil
}

// validateValue validates the metric value.
func validateValue(value float64) error {
	if math.IsNaN(value) || math.IsInf(value, 0) {
		return &validationError{
			errorType: "validation_failed",
			message:   "Invalid value: must be a valid number",
		}
	}
	return nil
}

// parseTimestamp parses an ISO 8601 timestamp or returns current time if empty.
func parseTimestamp(ts string) (time.Time, error) {
	if ts == "" {
		return time.Now().UTC(), nil
	}

	parsed, err := time.Parse(time.RFC3339, ts)
	if err != nil {
		return time.Time{}, &validationError{
			errorType: "validation_failed",
			message:   "Invalid timestamp format: must be ISO 8601 (e.g., 2024-01-15T10:30:00Z)",
		}
	}
	return parsed, nil
}

// validateMetadata validates the metadata map.
func validateMetadata(metadata map[string]string) error {
	if metadata == nil {
		return nil
	}

	if len(metadata) > MaxMetadataKeys {
		return &validationError{
			errorType: "validation_failed",
			message:   fmt.Sprintf("Metadata exceeds maximum of %d keys", MaxMetadataKeys),
		}
	}

	for key, value := range metadata {
		if key == "" {
			return &validationError{
				errorType: "validation_failed",
				message:   "Metadata key cannot be empty",
			}
		}
		if len(key) > MaxMetadataKeyLength {
			return &validationError{
				errorType: "validation_failed",
				message:   fmt.Sprintf("Metadata key '%s' exceeds maximum length of %d characters", key, MaxMetadataKeyLength),
			}
		}
		if len(value) > MaxMetadataValueLength {
			return &validationError{
				errorType: "validation_failed",
				message:   fmt.Sprintf("Metadata value for key '%s' exceeds maximum length of %d characters", key, MaxMetadataValueLength),
			}
		}
	}

	return nil
}

// GetMeasurementNames retrieves distinct measurement names with their metadata keys for a data source.
func (s *Service) GetMeasurementNames(ctx context.Context, dataSourceID uuid.UUID) ([]MeasurementSummary, error) {
	return s.repo.GetMeasurementNames(ctx, dataSourceID)
}

// GetMetadataValues retrieves all unique metadata key-value combinations for a specific measurement.
func (s *Service) GetMetadataValues(ctx context.Context, dataSourceID uuid.UUID, measurementName string) ([]MetadataValues, error) {
	return s.repo.GetMetadataValues(ctx, dataSourceID, measurementName)
}

// GetAggregatedMeasurements retrieves daily aggregated values with optional metadata filtering.
func (s *Service) GetAggregatedMeasurements(ctx context.Context, dataSourceID uuid.UUID, name string, startDate, endDate time.Time, metadataFilters map[string]string) ([]AggregatedDataPoint, error) {
	return s.repo.GetAggregatedMeasurements(ctx, dataSourceID, name, startDate, endDate, metadataFilters)
}

// GetRawMeasurements retrieves raw measurement data points with optional metadata filtering.
func (s *Service) GetRawMeasurements(ctx context.Context, dataSourceID uuid.UUID, name string, startDate, endDate time.Time, metadataFilters map[string]string, limit int) ([]Measurement, error) {
	return s.repo.GetRawMeasurements(ctx, dataSourceID, name, startDate, endDate, metadataFilters, limit)
}

// maxSplitSeries is the maximum number of distinct series to return before grouping into "Other".
const maxSplitSeries = 10

// GetAggregatedMeasurementsSplitBy retrieves daily aggregated values split by a metadata key.
// Keeps top 10 series by total sum and aggregates the rest into "Other".
func (s *Service) GetAggregatedMeasurementsSplitBy(ctx context.Context, dataSourceID uuid.UUID, name string, startDate, endDate time.Time, metadataFilters map[string]string, splitByKey string) ([]SplitSeries, error) {
	// Get raw split data from repository
	series, err := s.repo.GetAggregatedMeasurementsSplitBy(ctx, dataSourceID, name, startDate, endDate, metadataFilters, splitByKey)
	if err != nil {
		return nil, err
	}

	// If 10 or fewer series, return as-is (sorted by total)
	if len(series) <= maxSplitSeries {
		return sortSeriesByTotal(series), nil
	}

	// Calculate total sum for each series
	type seriesWithTotal struct {
		series SplitSeries
		total  float64
	}
	seriesWithTotals := make([]seriesWithTotal, len(series))
	for i, s := range series {
		var total float64
		for _, dp := range s.DataPoints {
			total += dp.Sum
		}
		seriesWithTotals[i] = seriesWithTotal{series: s, total: total}
	}

	// Sort by total descending
	sort.Slice(seriesWithTotals, func(i, j int) bool {
		return seriesWithTotals[i].total > seriesWithTotals[j].total
	})

	// Take top 10
	topSeries := make([]SplitSeries, maxSplitSeries)
	for i := 0; i < maxSplitSeries; i++ {
		topSeries[i] = seriesWithTotals[i].series
	}

	// Aggregate remaining into "Other"
	otherDataPoints := make(map[string]*AggregatedDataPoint)
	for i := maxSplitSeries; i < len(seriesWithTotals); i++ {
		for _, dp := range seriesWithTotals[i].series.DataPoints {
			if existing, ok := otherDataPoints[dp.Date]; ok {
				existing.Sum += dp.Sum
				existing.Count += dp.Count
			} else {
				otherDataPoints[dp.Date] = &AggregatedDataPoint{
					Date:  dp.Date,
					Sum:   dp.Sum,
					Count: dp.Count,
				}
			}
		}
	}

	// Convert "Other" map to sorted slice
	if len(otherDataPoints) > 0 {
		otherDPs := make([]AggregatedDataPoint, 0, len(otherDataPoints))
		for _, dp := range otherDataPoints {
			otherDPs = append(otherDPs, *dp)
		}
		sort.Slice(otherDPs, func(i, j int) bool {
			return otherDPs[i].Date < otherDPs[j].Date
		})
		topSeries = append(topSeries, SplitSeries{
			Key:        "Other",
			DataPoints: otherDPs,
		})
	}

	return topSeries, nil
}

// sortSeriesByTotal sorts series by their total sum in descending order.
func sortSeriesByTotal(series []SplitSeries) []SplitSeries {
	type seriesWithTotal struct {
		series SplitSeries
		total  float64
	}
	seriesWithTotals := make([]seriesWithTotal, len(series))
	for i, s := range series {
		var total float64
		for _, dp := range s.DataPoints {
			total += dp.Sum
		}
		seriesWithTotals[i] = seriesWithTotal{series: s, total: total}
	}

	sort.Slice(seriesWithTotals, func(i, j int) bool {
		return seriesWithTotals[i].total > seriesWithTotals[j].total
	})

	result := make([]SplitSeries, len(series))
	for i, s := range seriesWithTotals {
		result[i] = s.series
	}
	return result
}
