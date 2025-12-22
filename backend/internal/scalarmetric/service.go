package scalarmetric

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/devbydaniel/litekpi/internal/datasource"
	"github.com/devbydaniel/litekpi/internal/ingest"
)

// Service handles scalar metric business logic.
type Service struct {
	repo              *Repository
	ingestService     *ingest.Service
	dataSourceService *datasource.Service
}

// NewService creates a new scalar metric service.
func NewService(repo *Repository, ingestService *ingest.Service, dataSourceService *datasource.Service) *Service {
	return &Service{
		repo:              repo,
		ingestService:     ingestService,
		dataSourceService: dataSourceService,
	}
}

// Create creates a new scalar metric.
// The caller is responsible for verifying dashboard ownership.
func (s *Service) Create(ctx context.Context, orgID, dashboardID uuid.UUID, req CreateScalarMetricRequest) (*ScalarMetric, error) {
	if err := s.validateCreateRequest(ctx, orgID, req); err != nil {
		return nil, err
	}

	maxPos, err := s.repo.GetMaxPosition(ctx, dashboardID)
	if err != nil {
		return nil, fmt.Errorf("failed to get max position: %w", err)
	}

	sm, err := s.repo.Create(ctx, dashboardID, req.DataSourceID, req, maxPos+1)
	if err != nil {
		return nil, fmt.Errorf("failed to create scalar metric: %w", err)
	}

	return sm, nil
}

func (s *Service) validateCreateRequest(ctx context.Context, orgID uuid.UUID, req CreateScalarMetricRequest) error {
	// Validate label
	label := strings.TrimSpace(req.Label)
	if label == "" {
		return ErrLabelEmpty
	}
	if len(label) > 255 {
		return ErrLabelTooLong
	}

	// Validate measurement name
	if strings.TrimSpace(req.MeasurementName) == "" {
		return ErrMeasurementNameEmpty
	}

	// Validate timeframe
	if !IsValidTimeframe(req.Timeframe) {
		return ErrInvalidTimeframe
	}

	// Validate aggregation
	if !req.Aggregation.IsValid() {
		return ErrInvalidAggregation
	}

	// Validate comparison display type if comparison is enabled
	if req.ComparisonEnabled && req.ComparisonDisplayType != nil {
		if !req.ComparisonDisplayType.IsValid() {
			return ErrInvalidComparisonType
		}
	}

	// Verify data source ownership
	_, err := s.dataSourceService.GetDataSource(ctx, orgID, req.DataSourceID)
	if err != nil {
		return fmt.Errorf("failed to verify data source: %w", err)
	}

	return nil
}

// GetByDashboardID retrieves all scalar metrics for a dashboard.
func (s *Service) GetByDashboardID(ctx context.Context, dashboardID uuid.UUID) ([]ScalarMetric, error) {
	scalarMetrics, err := s.repo.GetByDashboardID(ctx, dashboardID)
	if err != nil {
		return nil, fmt.Errorf("failed to get scalar metrics: %w", err)
	}
	return scalarMetrics, nil
}

// GetByID retrieves a scalar metric by its ID.
func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*ScalarMetric, error) {
	sm, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get scalar metric: %w", err)
	}
	if sm == nil {
		return nil, ErrScalarMetricNotFound
	}
	return sm, nil
}

// Update updates a scalar metric's configuration.
// The caller is responsible for verifying dashboard ownership.
func (s *Service) Update(ctx context.Context, dashboardID, scalarMetricID uuid.UUID, req UpdateScalarMetricRequest) (*ScalarMetric, error) {
	// Verify scalar metric exists and belongs to dashboard
	sm, err := s.repo.GetByID(ctx, scalarMetricID)
	if err != nil {
		return nil, fmt.Errorf("failed to get scalar metric: %w", err)
	}
	if sm == nil {
		return nil, ErrScalarMetricNotFound
	}
	if sm.DashboardID != dashboardID {
		return nil, ErrScalarMetricNotFound
	}

	// Validate label
	label := strings.TrimSpace(req.Label)
	if label == "" {
		return nil, ErrLabelEmpty
	}
	if len(label) > 255 {
		return nil, ErrLabelTooLong
	}

	// Validate timeframe
	if !IsValidTimeframe(req.Timeframe) {
		return nil, ErrInvalidTimeframe
	}

	// Validate aggregation
	if !req.Aggregation.IsValid() {
		return nil, ErrInvalidAggregation
	}

	// Validate comparison display type if comparison is enabled
	if req.ComparisonEnabled && req.ComparisonDisplayType != nil {
		if !req.ComparisonDisplayType.IsValid() {
			return nil, ErrInvalidComparisonType
		}
	}

	if err := s.repo.Update(ctx, scalarMetricID, req); err != nil {
		return nil, fmt.Errorf("failed to update scalar metric: %w", err)
	}

	// Return updated scalar metric
	return s.repo.GetByID(ctx, scalarMetricID)
}

// Delete deletes a scalar metric.
// The caller is responsible for verifying dashboard ownership.
func (s *Service) Delete(ctx context.Context, dashboardID, scalarMetricID uuid.UUID) error {
	// Verify scalar metric exists and belongs to dashboard
	sm, err := s.repo.GetByID(ctx, scalarMetricID)
	if err != nil {
		return fmt.Errorf("failed to get scalar metric: %w", err)
	}
	if sm == nil {
		return ErrScalarMetricNotFound
	}
	if sm.DashboardID != dashboardID {
		return ErrScalarMetricNotFound
	}

	if err := s.repo.Delete(ctx, scalarMetricID); err != nil {
		return fmt.Errorf("failed to delete scalar metric: %w", err)
	}
	return nil
}

// Reorder reorders scalar metrics on a dashboard.
// The caller is responsible for verifying dashboard ownership.
func (s *Service) Reorder(ctx context.Context, dashboardID uuid.UUID, scalarMetricIDs []uuid.UUID) error {
	if err := s.repo.UpdatePositions(ctx, dashboardID, scalarMetricIDs); err != nil {
		return fmt.Errorf("failed to reorder scalar metrics: %w", err)
	}
	return nil
}

// Compute calculates the values for a list of scalar metrics.
func (s *Service) Compute(ctx context.Context, scalarMetrics []ScalarMetric) ([]ComputedScalarMetric, error) {
	computed := make([]ComputedScalarMetric, len(scalarMetrics))

	for i, sm := range scalarMetrics {
		result, err := s.computeOne(ctx, sm)
		if err != nil {
			return nil, fmt.Errorf("failed to compute scalar metric %s: %w", sm.ID, err)
		}
		computed[i] = *result
	}

	return computed, nil
}

func (s *Service) computeOne(ctx context.Context, sm ScalarMetric) (*ComputedScalarMetric, error) {
	// Calculate date ranges
	currentStart, currentEnd := getTimeframeRange(sm.Timeframe)

	// Build metadata filters
	filters := make(map[string]string)
	for _, f := range sm.Filters {
		filters[f.Key] = f.Value
	}

	// Query current period
	currentData, err := s.ingestService.GetAggregatedMeasurements(
		ctx, sm.DataSourceID, sm.MeasurementName,
		currentStart, currentEnd, filters,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get current period data: %w", err)
	}

	// Calculate value based on aggregation
	value := aggregate(currentData, sm.Aggregation)

	computed := &ComputedScalarMetric{
		ScalarMetric: sm,
		Value:        value,
	}

	// Handle comparison if enabled
	if sm.ComparisonEnabled {
		previousStart, previousEnd := getPreviousTimeframeRange(sm.Timeframe, currentStart, currentEnd)

		previousData, err := s.ingestService.GetAggregatedMeasurements(
			ctx, sm.DataSourceID, sm.MeasurementName,
			previousStart, previousEnd, filters,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to get previous period data: %w", err)
		}

		previousValue := aggregate(previousData, sm.Aggregation)
		computed.PreviousValue = &previousValue

		change := value - previousValue
		computed.Change = &change

		if previousValue != 0 {
			changePercent := (change / previousValue) * 100
			computed.ChangePercent = &changePercent
		}
	}

	return computed, nil
}

func aggregate(data []ingest.AggregatedDataPoint, aggregationType Aggregation) float64 {
	if len(data) == 0 {
		return 0
	}

	var totalSum float64
	var totalCount int
	for _, dp := range data {
		totalSum += dp.Sum
		totalCount += dp.Count
	}

	switch aggregationType {
	case AggregationAverage:
		if totalCount == 0 {
			return 0
		}
		return totalSum / float64(totalCount)
	default: // AggregationSum
		return totalSum
	}
}

func getTimeframeRange(timeframe string) (start, end time.Time) {
	now := time.Now().UTC()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	switch timeframe {
	case "last_7_days":
		start = today.AddDate(0, 0, -7)
		end = today.AddDate(0, 0, 1) // Include today
	case "last_30_days":
		start = today.AddDate(0, 0, -30)
		end = today.AddDate(0, 0, 1)
	case "this_month":
		start = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
		end = today.AddDate(0, 0, 1)
	case "last_month":
		firstOfThisMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
		start = firstOfThisMonth.AddDate(0, -1, 0)
		end = firstOfThisMonth
	default:
		// Default to last 30 days
		start = today.AddDate(0, 0, -30)
		end = today.AddDate(0, 0, 1)
	}

	return start, end
}

func getPreviousTimeframeRange(timeframe string, currentStart, currentEnd time.Time) (start, end time.Time) {
	duration := currentEnd.Sub(currentStart)

	switch timeframe {
	case "last_7_days":
		// 7 days before the current 7-day window
		end = currentStart
		start = end.Add(-duration)
	case "last_30_days":
		// 30 days before the current 30-day window
		end = currentStart
		start = end.Add(-duration)
	case "this_month":
		// Same date range in the previous month
		start = currentStart.AddDate(0, -1, 0)
		end = currentEnd.AddDate(0, -1, 0)
	case "last_month":
		// The month before last month
		start = currentStart.AddDate(0, -1, 0)
		end = currentEnd.AddDate(0, -1, 0)
	default:
		// Default: same duration before
		end = currentStart
		start = end.Add(-duration)
	}

	return start, end
}
