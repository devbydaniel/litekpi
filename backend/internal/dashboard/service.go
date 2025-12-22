package dashboard

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"

	"github.com/devbydaniel/litekpi/internal/scalarmetric"
	"github.com/devbydaniel/litekpi/internal/timeseries"
)

// Service handles dashboard business logic.
type Service struct {
	repo                *Repository
	timeSeriesService   *timeseries.Service
	scalarMetricService *scalarmetric.Service
}

// NewService creates a new dashboard service.
func NewService(repo *Repository, timeSeriesService *timeseries.Service, scalarMetricService *scalarmetric.Service) *Service {
	return &Service{
		repo:                repo,
		timeSeriesService:   timeSeriesService,
		scalarMetricService: scalarMetricService,
	}
}

// CreateDashboard creates a new dashboard.
func (s *Service) CreateDashboard(ctx context.Context, orgID uuid.UUID, req CreateDashboardRequest) (*Dashboard, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, ErrDashboardNameEmpty
	}

	dashboard, err := s.repo.CreateDashboard(ctx, orgID, name, false)
	if err != nil {
		return nil, fmt.Errorf("failed to create dashboard: %w", err)
	}

	return dashboard, nil
}

// ListDashboards returns all dashboards for an organization.
func (s *Service) ListDashboards(ctx context.Context, orgID uuid.UUID) ([]Dashboard, error) {
	dashboards, err := s.repo.GetDashboardsByOrganizationID(ctx, orgID)
	if err != nil {
		return nil, fmt.Errorf("failed to list dashboards: %w", err)
	}

	return dashboards, nil
}

// GetDashboard returns a dashboard with its time series and scalar metrics after verifying organization ownership.
func (s *Service) GetDashboard(ctx context.Context, orgID, dashboardID uuid.UUID) (*DashboardWithData, error) {
	dashboard, err := s.repo.GetDashboardByID(ctx, dashboardID)
	if err != nil {
		return nil, fmt.Errorf("failed to get dashboard: %w", err)
	}
	if dashboard == nil {
		return nil, ErrDashboardNotFound
	}
	if dashboard.OrganizationID != orgID {
		return nil, ErrUnauthorized
	}

	ts, err := s.timeSeriesService.GetByDashboardID(ctx, dashboardID)
	if err != nil {
		return nil, fmt.Errorf("failed to get time series: %w", err)
	}

	sm, err := s.scalarMetricService.GetByDashboardID(ctx, dashboardID)
	if err != nil {
		return nil, fmt.Errorf("failed to get scalar metrics: %w", err)
	}

	return &DashboardWithData{
		Dashboard:     *dashboard,
		TimeSeries:    ts,
		ScalarMetrics: sm,
	}, nil
}

// GetDefaultDashboard returns the default dashboard for an organization.
func (s *Service) GetDefaultDashboard(ctx context.Context, orgID uuid.UUID) (*DashboardWithData, error) {
	dashboard, err := s.repo.GetDefaultDashboard(ctx, orgID)
	if err != nil {
		return nil, fmt.Errorf("failed to get default dashboard: %w", err)
	}
	if dashboard == nil {
		// Create default dashboard if it doesn't exist
		dashboard, err = s.repo.CreateDashboard(ctx, orgID, "Dashboard", true)
		if err != nil {
			return nil, fmt.Errorf("failed to create default dashboard: %w", err)
		}
	}

	ts, err := s.timeSeriesService.GetByDashboardID(ctx, dashboard.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get time series: %w", err)
	}

	sm, err := s.scalarMetricService.GetByDashboardID(ctx, dashboard.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get scalar metrics: %w", err)
	}

	return &DashboardWithData{
		Dashboard:     *dashboard,
		TimeSeries:    ts,
		ScalarMetrics: sm,
	}, nil
}

// UpdateDashboard updates a dashboard's name.
func (s *Service) UpdateDashboard(ctx context.Context, orgID, dashboardID uuid.UUID, req UpdateDashboardRequest) (*Dashboard, error) {
	dashboard, err := s.repo.GetDashboardByID(ctx, dashboardID)
	if err != nil {
		return nil, fmt.Errorf("failed to get dashboard: %w", err)
	}
	if dashboard == nil {
		return nil, ErrDashboardNotFound
	}
	if dashboard.OrganizationID != orgID {
		return nil, ErrUnauthorized
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, ErrDashboardNameEmpty
	}

	if err := s.repo.UpdateDashboard(ctx, dashboardID, name); err != nil {
		return nil, fmt.Errorf("failed to update dashboard: %w", err)
	}

	dashboard.Name = name
	return dashboard, nil
}

// DeleteDashboard deletes a dashboard.
func (s *Service) DeleteDashboard(ctx context.Context, orgID, dashboardID uuid.UUID) error {
	dashboard, err := s.repo.GetDashboardByID(ctx, dashboardID)
	if err != nil {
		return fmt.Errorf("failed to get dashboard: %w", err)
	}
	if dashboard == nil {
		return ErrDashboardNotFound
	}
	if dashboard.OrganizationID != orgID {
		return ErrUnauthorized
	}
	if dashboard.IsDefault {
		return ErrCannotDeleteDefault
	}

	if err := s.repo.DeleteDashboard(ctx, dashboardID); err != nil {
		return fmt.Errorf("failed to delete dashboard: %w", err)
	}

	return nil
}

// VerifyDashboardOwnership verifies that a dashboard belongs to an organization.
// This is used by handlers to check ownership before delegating to time series or scalar metric services.
func (s *Service) VerifyDashboardOwnership(ctx context.Context, orgID, dashboardID uuid.UUID) (*Dashboard, error) {
	dashboard, err := s.repo.GetDashboardByID(ctx, dashboardID)
	if err != nil {
		return nil, fmt.Errorf("failed to get dashboard: %w", err)
	}
	if dashboard == nil {
		return nil, ErrDashboardNotFound
	}
	if dashboard.OrganizationID != orgID {
		return nil, ErrUnauthorized
	}
	return dashboard, nil
}
