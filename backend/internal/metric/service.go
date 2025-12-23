package metric

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/devbydaniel/litekpi/internal/datasource"
)

const maxSplitBySeries = 10 // Maximum number of series when using split_by

// Service handles metric business logic.
type Service struct {
	repo              *Repository
	dataSourceService *datasource.Service
}

// NewService creates a new metric service.
func NewService(repo *Repository, dataSourceService *datasource.Service) *Service {
	return &Service{
		repo:              repo,
		dataSourceService: dataSourceService,
	}
}

// Create creates a new metric.
// The caller is responsible for verifying dashboard ownership.
func (s *Service) Create(ctx context.Context, orgID, dashboardID uuid.UUID, req CreateMetricRequest) (*Metric, error) {
	if err := s.validateCreateRequest(ctx, orgID, req); err != nil {
		return nil, err
	}

	maxPos, err := s.repo.GetMaxPosition(ctx, dashboardID)
	if err != nil {
		return nil, fmt.Errorf("failed to get max position: %w", err)
	}

	m, err := s.repo.Create(ctx, dashboardID, req.DataSourceID, req, maxPos+1)
	if err != nil {
		return nil, fmt.Errorf("failed to create metric: %w", err)
	}

	return m, nil
}

func (s *Service) validateCreateRequest(ctx context.Context, orgID uuid.UUID, req CreateMetricRequest) error {
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

	// Validate aggregation_key for count_unique
	if req.Aggregation.RequiresAggregationKey() {
		if req.AggregationKey == nil || strings.TrimSpace(*req.AggregationKey) == "" {
			return ErrAggregationKeyRequired
		}
	}

	// Validate display mode
	if !req.DisplayMode.IsValid() {
		return ErrInvalidDisplayMode
	}

	// Validate display mode specific fields
	if req.DisplayMode == DisplayModeTimeSeries {
		// Granularity is required for time series
		if req.Granularity == nil || !req.Granularity.IsValid() {
			return ErrInvalidGranularity
		}
		if req.ChartType == nil || !req.ChartType.IsValid() {
			return ErrChartTypeRequired
		}
	}
	// For scalar, granularity should be nil (ignored if provided)

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

// GetByDashboardID retrieves all metrics for a dashboard.
func (s *Service) GetByDashboardID(ctx context.Context, dashboardID uuid.UUID) ([]Metric, error) {
	metrics, err := s.repo.GetByDashboardID(ctx, dashboardID)
	if err != nil {
		return nil, fmt.Errorf("failed to get metrics: %w", err)
	}
	return metrics, nil
}

// GetByID retrieves a metric by its ID.
func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*Metric, error) {
	m, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get metric: %w", err)
	}
	if m == nil {
		return nil, ErrMetricNotFound
	}
	return m, nil
}

// Update updates a metric's configuration.
// The caller is responsible for verifying dashboard ownership.
func (s *Service) Update(ctx context.Context, dashboardID, metricID uuid.UUID, req UpdateMetricRequest) (*Metric, error) {
	// Verify metric exists and belongs to dashboard
	m, err := s.repo.GetByID(ctx, metricID)
	if err != nil {
		return nil, fmt.Errorf("failed to get metric: %w", err)
	}
	if m == nil {
		return nil, ErrMetricNotFound
	}
	if m.DashboardID != dashboardID {
		return nil, ErrMetricNotFound
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

	// Validate aggregation_key for count_unique
	if req.Aggregation.RequiresAggregationKey() {
		if req.AggregationKey == nil || strings.TrimSpace(*req.AggregationKey) == "" {
			return nil, ErrAggregationKeyRequired
		}
	}

	// Validate display mode
	if !req.DisplayMode.IsValid() {
		return nil, ErrInvalidDisplayMode
	}

	// Validate display mode specific fields
	if req.DisplayMode == DisplayModeTimeSeries {
		// Granularity is required for time series
		if req.Granularity == nil || !req.Granularity.IsValid() {
			return nil, ErrInvalidGranularity
		}
		if req.ChartType == nil || !req.ChartType.IsValid() {
			return nil, ErrChartTypeRequired
		}
	}
	// For scalar, granularity should be nil (ignored if provided)

	// Validate comparison display type if comparison is enabled
	if req.ComparisonEnabled && req.ComparisonDisplayType != nil {
		if !req.ComparisonDisplayType.IsValid() {
			return nil, ErrInvalidComparisonType
		}
	}

	if err := s.repo.Update(ctx, metricID, req); err != nil {
		return nil, fmt.Errorf("failed to update metric: %w", err)
	}

	// Return updated metric
	return s.repo.GetByID(ctx, metricID)
}

// Delete deletes a metric.
// The caller is responsible for verifying dashboard ownership.
func (s *Service) Delete(ctx context.Context, dashboardID, metricID uuid.UUID) error {
	// Verify metric exists and belongs to dashboard
	m, err := s.repo.GetByID(ctx, metricID)
	if err != nil {
		return fmt.Errorf("failed to get metric: %w", err)
	}
	if m == nil {
		return ErrMetricNotFound
	}
	if m.DashboardID != dashboardID {
		return ErrMetricNotFound
	}

	if err := s.repo.Delete(ctx, metricID); err != nil {
		return fmt.Errorf("failed to delete metric: %w", err)
	}
	return nil
}

// Reorder reorders metrics on a dashboard.
// The caller is responsible for verifying dashboard ownership.
func (s *Service) Reorder(ctx context.Context, dashboardID uuid.UUID, metricIDs []uuid.UUID) error {
	if err := s.repo.UpdatePositions(ctx, dashboardID, metricIDs); err != nil {
		return fmt.Errorf("failed to reorder metrics: %w", err)
	}
	return nil
}

// Compute calculates the values for a list of metrics.
func (s *Service) Compute(ctx context.Context, metrics []Metric) ([]ComputedMetric, error) {
	computed := make([]ComputedMetric, len(metrics))

	for i, m := range metrics {
		result, err := s.computeOne(ctx, m)
		if err != nil {
			return nil, fmt.Errorf("failed to compute metric %s: %w", m.ID, err)
		}
		computed[i] = *result
	}

	return computed, nil
}

func (s *Service) computeOne(ctx context.Context, m Metric) (*ComputedMetric, error) {
	// Calculate date ranges
	currentStart, currentEnd := getTimeframeRange(m.Timeframe, m.DateFrom, m.DateTo)

	// Build metadata filters
	filters := make(map[string]string)
	for _, f := range m.Filters {
		filters[f.Key] = f.Value
	}

	computed := &ComputedMetric{Metric: m}

	switch m.DisplayMode {
	case DisplayModeScalar:
		return s.computeScalar(ctx, m, currentStart, currentEnd, filters)
	case DisplayModeTimeSeries:
		return s.computeTimeSeries(ctx, m, currentStart, currentEnd, filters)
	}

	return computed, nil
}

func (s *Service) computeScalar(ctx context.Context, m Metric, start, end time.Time, filters map[string]string) (*ComputedMetric, error) {
	computed := &ComputedMetric{Metric: m}

	// Get current value
	value, err := s.aggregateScalarValue(ctx, m, start, end, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to get current period data: %w", err)
	}
	computed.Value = &value

	// Handle comparison if enabled
	if m.ComparisonEnabled {
		previousStart, previousEnd := getPreviousTimeframeRange(m.Timeframe, start, end)

		previousValue, err := s.aggregateScalarValue(ctx, m, previousStart, previousEnd, filters)
		if err != nil {
			return nil, fmt.Errorf("failed to get previous period data: %w", err)
		}
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

func (s *Service) aggregateScalarValue(ctx context.Context, m Metric, start, end time.Time, filters map[string]string) (float64, error) {
	switch m.Aggregation {
	case AggregationCountUnique:
		// Use scalar method - correctly counts unique values across entire timeframe
		count, err := s.repo.GetScalarCountUnique(ctx, m.DataSourceID, m.MeasurementName, start, end, filters, *m.AggregationKey)
		if err != nil {
			return 0, err
		}
		return float64(count), nil

	case AggregationCount:
		_, count, err := s.repo.GetScalarAggregate(ctx, m.DataSourceID, m.MeasurementName, start, end, filters)
		if err != nil {
			return 0, err
		}
		return float64(count), nil

	case AggregationAverage:
		sum, count, err := s.repo.GetScalarAggregate(ctx, m.DataSourceID, m.MeasurementName, start, end, filters)
		if err != nil {
			return 0, err
		}
		if count == 0 {
			return 0, nil
		}
		return sum / float64(count), nil

	default: // sum
		sum, _, err := s.repo.GetScalarAggregate(ctx, m.DataSourceID, m.MeasurementName, start, end, filters)
		if err != nil {
			return 0, err
		}
		return sum, nil
	}
}

func (s *Service) computeTimeSeries(ctx context.Context, m Metric, start, end time.Time, filters map[string]string) (*ComputedMetric, error) {
	if m.Granularity == nil {
		return nil, fmt.Errorf("granularity is required for time series metrics")
	}

	computed := &ComputedMetric{Metric: m}

	if m.SplitBy != nil && *m.SplitBy != "" {
		series, err := s.getTimeSeriesSplitBy(ctx, m, start, end, filters)
		if err != nil {
			return nil, fmt.Errorf("failed to get split time series data: %w", err)
		}
		computed.Series = series
	} else {
		dataPoints, err := s.getTimeSeriesData(ctx, m, start, end, filters)
		if err != nil {
			return nil, fmt.Errorf("failed to get time series data: %w", err)
		}
		computed.DataPoints = dataPoints
	}

	return computed, nil
}

func (s *Service) getTimeSeriesData(ctx context.Context, m Metric, start, end time.Time, filters map[string]string) ([]DataPoint, error) {
	granularity := *m.Granularity // Already validated in computeTimeSeries

	switch m.Aggregation {
	case AggregationCountUnique:
		data, err := s.repo.GetCountUniqueMeasurements(ctx, m.DataSourceID, m.MeasurementName, start, end, filters, *m.AggregationKey, granularity)
		if err != nil {
			return nil, err
		}
		dataPoints := make([]DataPoint, len(data))
		for i, dp := range data {
			dataPoints[i] = DataPoint{
				Date:  dp.Date,
				Value: dp.Sum, // Sum holds the unique count
			}
		}
		return dataPoints, nil

	case AggregationCount:
		data, err := s.repo.GetAggregatedMeasurements(ctx, m.DataSourceID, m.MeasurementName, start, end, filters, granularity)
		if err != nil {
			return nil, err
		}
		dataPoints := make([]DataPoint, len(data))
		for i, dp := range data {
			dataPoints[i] = DataPoint{
				Date:  dp.Date,
				Value: float64(dp.Count),
			}
		}
		return dataPoints, nil

	case AggregationAverage:
		data, err := s.repo.GetAggregatedMeasurements(ctx, m.DataSourceID, m.MeasurementName, start, end, filters, granularity)
		if err != nil {
			return nil, err
		}
		dataPoints := make([]DataPoint, len(data))
		for i, dp := range data {
			var value float64
			if dp.Count > 0 {
				value = dp.Sum / float64(dp.Count)
			}
			dataPoints[i] = DataPoint{
				Date:  dp.Date,
				Value: value,
			}
		}
		return dataPoints, nil

	default: // sum
		data, err := s.repo.GetAggregatedMeasurements(ctx, m.DataSourceID, m.MeasurementName, start, end, filters, granularity)
		if err != nil {
			return nil, err
		}
		dataPoints := make([]DataPoint, len(data))
		for i, dp := range data {
			dataPoints[i] = DataPoint{
				Date:  dp.Date,
				Value: dp.Sum,
			}
		}
		return dataPoints, nil
	}
}

func (s *Service) getTimeSeriesSplitBy(ctx context.Context, m Metric, start, end time.Time, filters map[string]string) ([]SplitSeries, error) {
	granularity := *m.Granularity // Already validated in computeTimeSeries

	// Note: count_unique with split_by would require different query logic
	// For now, we only support sum/average/count with split_by
	if m.Aggregation == AggregationCountUnique {
		// Fall back to non-split behavior for count_unique
		dataPoints, err := s.getTimeSeriesData(ctx, m, start, end, filters)
		if err != nil {
			return nil, err
		}
		return []SplitSeries{{Key: "total", DataPoints: dataPoints}}, nil
	}

	series, err := s.repo.GetAggregatedMeasurementsSplitBy(ctx, m.DataSourceID, m.MeasurementName, start, end, filters, *m.SplitBy, granularity)
	if err != nil {
		return nil, err
	}

	// Apply top-N aggregation
	return applyTopNSeries(series, maxSplitBySeries), nil
}

// Helper functions

func getTimeframeRange(timeframe string, dateFrom, dateTo *time.Time) (start, end time.Time) {
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
	case "custom":
		if dateFrom != nil && dateTo != nil {
			start = *dateFrom
			end = dateTo.AddDate(0, 0, 1) // Include end date
		} else {
			// Fallback to last 30 days if custom dates not provided
			start = today.AddDate(0, 0, -30)
			end = today.AddDate(0, 0, 1)
		}
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
		end = currentStart
		start = end.Add(-duration)
	case "last_30_days":
		end = currentStart
		start = end.Add(-duration)
	case "this_month":
		start = currentStart.AddDate(0, -1, 0)
		end = currentEnd.AddDate(0, -1, 0)
	case "last_month":
		start = currentStart.AddDate(0, -1, 0)
		end = currentEnd.AddDate(0, -1, 0)
	case "custom":
		end = currentStart
		start = end.Add(-duration)
	default:
		end = currentStart
		start = end.Add(-duration)
	}

	return start, end
}

func applyTopNSeries(series []SplitSeries, maxSeries int) []SplitSeries {
	if len(series) <= maxSeries {
		return series
	}

	// Calculate total for each series
	type seriesTotal struct {
		series SplitSeries
		total  float64
	}
	totals := make([]seriesTotal, len(series))
	for i, s := range series {
		var total float64
		for _, dp := range s.DataPoints {
			total += dp.Value
		}
		totals[i] = seriesTotal{series: s, total: total}
	}

	// Sort by total descending
	sort.Slice(totals, func(i, j int) bool {
		return totals[i].total > totals[j].total
	})

	// Take top N-1 and aggregate the rest into "Other"
	result := make([]SplitSeries, 0, maxSeries)
	for i := 0; i < maxSeries-1 && i < len(totals); i++ {
		result = append(result, totals[i].series)
	}

	// Aggregate remaining into "Other"
	if len(totals) > maxSeries-1 {
		otherDataPoints := make(map[string]float64)
		for i := maxSeries - 1; i < len(totals); i++ {
			for _, dp := range totals[i].series.DataPoints {
				otherDataPoints[dp.Date] += dp.Value
			}
		}

		// Convert map to sorted slice
		var otherDps []DataPoint
		for date, value := range otherDataPoints {
			otherDps = append(otherDps, DataPoint{Date: date, Value: value})
		}
		sort.Slice(otherDps, func(i, j int) bool {
			return otherDps[i].Date < otherDps[j].Date
		})

		result = append(result, SplitSeries{Key: "Other", DataPoints: otherDps})
	}

	return result
}
