package kpi

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/devbydaniel/litekpi/internal/datasource"
	"github.com/devbydaniel/litekpi/internal/ingest"
)

// Service handles KPI business logic.
type Service struct {
	repo              *Repository
	ingestService     *ingest.Service
	dataSourceService *datasource.Service
}

// NewService creates a new KPI service.
func NewService(repo *Repository, ingestService *ingest.Service, dataSourceService *datasource.Service) *Service {
	return &Service{
		repo:              repo,
		ingestService:     ingestService,
		dataSourceService: dataSourceService,
	}
}

// CreateKPIForDashboard creates a new KPI for a dashboard.
func (s *Service) CreateKPIForDashboard(ctx context.Context, orgID, dashboardID uuid.UUID, req CreateKPIRequest) (*KPI, error) {
	if err := s.validateCreateRequest(ctx, orgID, req); err != nil {
		return nil, err
	}

	maxPos, err := s.repo.GetMaxKPIPositionForDashboard(ctx, dashboardID)
	if err != nil {
		return nil, fmt.Errorf("failed to get max KPI position: %w", err)
	}

	kpi, err := s.repo.CreateKPIForDashboard(ctx, dashboardID, req.DataSourceID, req, maxPos+1)
	if err != nil {
		return nil, fmt.Errorf("failed to create KPI: %w", err)
	}

	return kpi, nil
}

// CreateKPIForReport creates a new KPI for a report.
func (s *Service) CreateKPIForReport(ctx context.Context, orgID, reportID uuid.UUID, req CreateKPIRequest) (*KPI, error) {
	if err := s.validateCreateRequest(ctx, orgID, req); err != nil {
		return nil, err
	}

	maxPos, err := s.repo.GetMaxKPIPositionForReport(ctx, reportID)
	if err != nil {
		return nil, fmt.Errorf("failed to get max KPI position: %w", err)
	}

	kpi, err := s.repo.CreateKPIForReport(ctx, reportID, req.DataSourceID, req, maxPos+1)
	if err != nil {
		return nil, fmt.Errorf("failed to create KPI: %w", err)
	}

	return kpi, nil
}

func (s *Service) validateCreateRequest(ctx context.Context, orgID uuid.UUID, req CreateKPIRequest) error {
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
	if !IsValidAggregation(req.Aggregation) {
		return ErrInvalidAggregation
	}

	// Validate comparison display type if comparison is enabled
	if req.ComparisonEnabled && req.ComparisonDisplayType != nil {
		if !IsValidComparisonDisplayType(*req.ComparisonDisplayType) {
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

// GetKPIsByDashboardID retrieves all KPIs for a dashboard.
func (s *Service) GetKPIsByDashboardID(ctx context.Context, dashboardID uuid.UUID) ([]KPI, error) {
	kpis, err := s.repo.GetKPIsByDashboardID(ctx, dashboardID)
	if err != nil {
		return nil, fmt.Errorf("failed to get KPIs: %w", err)
	}
	return kpis, nil
}

// GetKPIsByReportID retrieves all KPIs for a report.
func (s *Service) GetKPIsByReportID(ctx context.Context, reportID uuid.UUID) ([]KPI, error) {
	kpis, err := s.repo.GetKPIsByReportID(ctx, reportID)
	if err != nil {
		return nil, fmt.Errorf("failed to get KPIs: %w", err)
	}
	return kpis, nil
}

// GetKPIByID retrieves a KPI by its ID.
func (s *Service) GetKPIByID(ctx context.Context, id uuid.UUID) (*KPI, error) {
	kpi, err := s.repo.GetKPIByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get KPI: %w", err)
	}
	if kpi == nil {
		return nil, ErrKPINotFound
	}
	return kpi, nil
}

// UpdateKPI updates a KPI's configuration.
func (s *Service) UpdateKPI(ctx context.Context, kpiID uuid.UUID, req UpdateKPIRequest) (*KPI, error) {
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
	if !IsValidAggregation(req.Aggregation) {
		return nil, ErrInvalidAggregation
	}

	// Validate comparison display type if comparison is enabled
	if req.ComparisonEnabled && req.ComparisonDisplayType != nil {
		if !IsValidComparisonDisplayType(*req.ComparisonDisplayType) {
			return nil, ErrInvalidComparisonType
		}
	}

	if err := s.repo.UpdateKPI(ctx, kpiID, req); err != nil {
		return nil, fmt.Errorf("failed to update KPI: %w", err)
	}

	// Return updated KPI
	return s.repo.GetKPIByID(ctx, kpiID)
}

// DeleteKPI deletes a KPI.
func (s *Service) DeleteKPI(ctx context.Context, kpiID uuid.UUID) error {
	if err := s.repo.DeleteKPI(ctx, kpiID); err != nil {
		return fmt.Errorf("failed to delete KPI: %w", err)
	}
	return nil
}

// ReorderKPIsForDashboard reorders KPIs on a dashboard.
func (s *Service) ReorderKPIsForDashboard(ctx context.Context, dashboardID uuid.UUID, kpiIDs []uuid.UUID) error {
	if err := s.repo.UpdateKPIPositionsForDashboard(ctx, dashboardID, kpiIDs); err != nil {
		return fmt.Errorf("failed to reorder KPIs: %w", err)
	}
	return nil
}

// ReorderKPIsForReport reorders KPIs on a report.
func (s *Service) ReorderKPIsForReport(ctx context.Context, reportID uuid.UUID, kpiIDs []uuid.UUID) error {
	if err := s.repo.UpdateKPIPositionsForReport(ctx, reportID, kpiIDs); err != nil {
		return fmt.Errorf("failed to reorder KPIs: %w", err)
	}
	return nil
}

// ComputeKPIs calculates the values for a list of KPIs.
func (s *Service) ComputeKPIs(ctx context.Context, kpis []KPI) ([]ComputedKPI, error) {
	computedKPIs := make([]ComputedKPI, len(kpis))

	for i, kpi := range kpis {
		computed, err := s.computeKPI(ctx, kpi)
		if err != nil {
			return nil, fmt.Errorf("failed to compute KPI %s: %w", kpi.ID, err)
		}
		computedKPIs[i] = *computed
	}

	return computedKPIs, nil
}

func (s *Service) computeKPI(ctx context.Context, kpi KPI) (*ComputedKPI, error) {
	// Calculate date ranges
	currentStart, currentEnd := getTimeframeRange(kpi.Timeframe)

	// Build metadata filters
	filters := make(map[string]string)
	for _, f := range kpi.Filters {
		filters[f.Key] = f.Value
	}

	// Query current period
	currentData, err := s.ingestService.GetAggregatedMeasurements(
		ctx, kpi.DataSourceID, kpi.MeasurementName,
		currentStart, currentEnd, filters,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get current period data: %w", err)
	}

	// Calculate value based on aggregation
	value := aggregate(currentData, kpi.Aggregation)

	computed := &ComputedKPI{
		KPI:   kpi,
		Value: value,
	}

	// Handle comparison if enabled
	if kpi.ComparisonEnabled {
		previousStart, previousEnd := getPreviousTimeframeRange(kpi.Timeframe, currentStart, currentEnd)

		previousData, err := s.ingestService.GetAggregatedMeasurements(
			ctx, kpi.DataSourceID, kpi.MeasurementName,
			previousStart, previousEnd, filters,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to get previous period data: %w", err)
		}

		previousValue := aggregate(previousData, kpi.Aggregation)
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

func aggregate(data []ingest.AggregatedDataPoint, aggregationType string) float64 {
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
	case "average":
		if totalCount == 0 {
			return 0
		}
		return totalSum / float64(totalCount)
	default: // "sum"
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
