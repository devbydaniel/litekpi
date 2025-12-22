package timeseries

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/devbydaniel/litekpi/internal/datasource"
)

// Service handles time series business logic.
type Service struct {
	repo              *Repository
	dataSourceService *datasource.Service
}

// NewService creates a new time series service.
func NewService(repo *Repository, dataSourceService *datasource.Service) *Service {
	return &Service{
		repo:              repo,
		dataSourceService: dataSourceService,
	}
}

// Create creates a new time series.
// The caller is responsible for verifying dashboard ownership.
func (s *Service) Create(ctx context.Context, orgID, dashboardID uuid.UUID, req CreateTimeSeriesRequest) (*TimeSeries, error) {
	// Verify data source ownership
	_, err := s.dataSourceService.GetDataSource(ctx, orgID, req.DataSourceID)
	if err != nil {
		return nil, fmt.Errorf("failed to verify data source: %w", err)
	}

	// Validate title length
	if req.Title != nil && len(*req.Title) > 128 {
		return nil, ErrTitleTooLong
	}

	// Validate measurement name
	if req.MeasurementName == "" {
		return nil, ErrMeasurementNameEmpty
	}

	// Validate chart type
	if !req.ChartType.IsValid() {
		return nil, ErrInvalidChartType
	}

	// Validate date range
	if !IsValidDateRange(req.DateRange) {
		return nil, ErrInvalidDateRange
	}

	// Get max position
	maxPos, err := s.repo.GetMaxPosition(ctx, dashboardID)
	if err != nil {
		return nil, fmt.Errorf("failed to get max position: %w", err)
	}

	ts, err := s.repo.Create(ctx, dashboardID, req.DataSourceID, req, maxPos+1)
	if err != nil {
		return nil, fmt.Errorf("failed to create time series: %w", err)
	}

	return ts, nil
}

// GetByDashboardID returns all time series for a dashboard.
func (s *Service) GetByDashboardID(ctx context.Context, dashboardID uuid.UUID) ([]TimeSeries, error) {
	timeSeries, err := s.repo.GetByDashboardID(ctx, dashboardID)
	if err != nil {
		return nil, fmt.Errorf("failed to get time series: %w", err)
	}

	return timeSeries, nil
}

// GetByID returns a time series by ID.
func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*TimeSeries, error) {
	ts, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get time series: %w", err)
	}
	if ts == nil {
		return nil, ErrTimeSeriesNotFound
	}

	return ts, nil
}

// Update updates a time series configuration.
// The caller is responsible for verifying dashboard ownership.
func (s *Service) Update(ctx context.Context, dashboardID, timeSeriesID uuid.UUID, req UpdateTimeSeriesRequest) (*TimeSeries, error) {
	// Verify time series exists and belongs to dashboard
	ts, err := s.repo.GetByID(ctx, timeSeriesID)
	if err != nil {
		return nil, fmt.Errorf("failed to get time series: %w", err)
	}
	if ts == nil {
		return nil, ErrTimeSeriesNotFound
	}
	if ts.DashboardID != dashboardID {
		return nil, ErrTimeSeriesNotFound
	}

	// Validate title length
	if req.Title != nil && len(*req.Title) > 128 {
		return nil, ErrTitleTooLong
	}

	// Validate chart type
	if !req.ChartType.IsValid() {
		return nil, ErrInvalidChartType
	}

	// Validate date range
	if !IsValidDateRange(req.DateRange) {
		return nil, ErrInvalidDateRange
	}

	if err := s.repo.Update(ctx, timeSeriesID, req); err != nil {
		return nil, fmt.Errorf("failed to update time series: %w", err)
	}

	// Return updated time series
	ts.Title = req.Title
	ts.ChartType = req.ChartType
	ts.DateRange = req.DateRange
	ts.DateFrom = req.DateFrom
	ts.DateTo = req.DateTo
	ts.SplitBy = req.SplitBy
	ts.Filters = req.Filters
	if ts.Filters == nil {
		ts.Filters = []Filter{}
	}

	return ts, nil
}

// Delete deletes a time series.
// The caller is responsible for verifying dashboard ownership.
func (s *Service) Delete(ctx context.Context, dashboardID, timeSeriesID uuid.UUID) error {
	// Verify time series exists and belongs to dashboard
	ts, err := s.repo.GetByID(ctx, timeSeriesID)
	if err != nil {
		return fmt.Errorf("failed to get time series: %w", err)
	}
	if ts == nil {
		return ErrTimeSeriesNotFound
	}
	if ts.DashboardID != dashboardID {
		return ErrTimeSeriesNotFound
	}

	if err := s.repo.Delete(ctx, timeSeriesID); err != nil {
		return fmt.Errorf("failed to delete time series: %w", err)
	}

	return nil
}

// Reorder reorders time series on a dashboard.
// The caller is responsible for verifying dashboard ownership.
func (s *Service) Reorder(ctx context.Context, dashboardID uuid.UUID, timeSeriesIDs []uuid.UUID) error {
	if err := s.repo.UpdatePositions(ctx, dashboardID, timeSeriesIDs); err != nil {
		return fmt.Errorf("failed to reorder time series: %w", err)
	}

	return nil
}
