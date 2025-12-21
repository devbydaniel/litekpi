package dashboard

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"

	"github.com/devbydaniel/litekpi/internal/datasource"
)

// Service handles dashboard business logic.
type Service struct {
	repo              *Repository
	dataSourceService *datasource.Service
}

// NewService creates a new dashboard service.
func NewService(repo *Repository, dataSourceService *datasource.Service) *Service {
	return &Service{
		repo:              repo,
		dataSourceService: dataSourceService,
	}
}

// Dashboard operations

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

// GetDashboard returns a dashboard with its widgets after verifying organization ownership.
func (s *Service) GetDashboard(ctx context.Context, orgID, dashboardID uuid.UUID) (*DashboardWithWidgets, error) {
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

	widgets, err := s.repo.GetWidgetsByDashboardID(ctx, dashboardID)
	if err != nil {
		return nil, fmt.Errorf("failed to get widgets: %w", err)
	}

	return &DashboardWithWidgets{
		Dashboard: *dashboard,
		Widgets:   widgets,
	}, nil
}

// GetDefaultDashboard returns the default dashboard for an organization.
func (s *Service) GetDefaultDashboard(ctx context.Context, orgID uuid.UUID) (*DashboardWithWidgets, error) {
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

	widgets, err := s.repo.GetWidgetsByDashboardID(ctx, dashboard.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get widgets: %w", err)
	}

	return &DashboardWithWidgets{
		Dashboard: *dashboard,
		Widgets:   widgets,
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

// Widget operations

// CreateWidget creates a new widget on a dashboard.
func (s *Service) CreateWidget(ctx context.Context, orgID, dashboardID uuid.UUID, req CreateWidgetRequest) (*Widget, error) {
	// Verify dashboard ownership
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

	// Verify data source ownership
	_, err = s.dataSourceService.GetDataSource(ctx, orgID, req.DataSourceID)
	if err != nil {
		return nil, fmt.Errorf("failed to verify data source: %w", err)
	}

	// Validate title length
	if req.Title != nil && len(*req.Title) > 128 {
		return nil, ErrTitleTooLong
	}

	// Validate chart type
	if !validChartTypes[req.ChartType] {
		return nil, ErrInvalidChartType
	}

	// Validate date range
	if !validDateRanges[req.DateRange] {
		return nil, ErrInvalidDateRange
	}

	// Get max position
	maxPos, err := s.repo.GetMaxWidgetPosition(ctx, dashboardID)
	if err != nil {
		return nil, fmt.Errorf("failed to get max widget position: %w", err)
	}

	widget, err := s.repo.CreateWidget(ctx, dashboardID, req.DataSourceID, req, maxPos+1)
	if err != nil {
		return nil, fmt.Errorf("failed to create widget: %w", err)
	}

	return widget, nil
}

// UpdateWidget updates a widget's configuration.
func (s *Service) UpdateWidget(ctx context.Context, orgID, dashboardID, widgetID uuid.UUID, req UpdateWidgetRequest) (*Widget, error) {
	// Verify dashboard ownership
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

	// Verify widget exists and belongs to dashboard
	widget, err := s.repo.GetWidgetByID(ctx, widgetID)
	if err != nil {
		return nil, fmt.Errorf("failed to get widget: %w", err)
	}
	if widget == nil {
		return nil, ErrWidgetNotFound
	}
	if widget.DashboardID != dashboardID {
		return nil, ErrWidgetNotFound
	}

	// Validate title length
	if req.Title != nil && len(*req.Title) > 128 {
		return nil, ErrTitleTooLong
	}

	// Validate chart type
	if !validChartTypes[req.ChartType] {
		return nil, ErrInvalidChartType
	}

	// Validate date range
	if !validDateRanges[req.DateRange] {
		return nil, ErrInvalidDateRange
	}

	if err := s.repo.UpdateWidget(ctx, widgetID, req); err != nil {
		return nil, fmt.Errorf("failed to update widget: %w", err)
	}

	// Return updated widget
	widget.Title = req.Title
	widget.ChartType = req.ChartType
	widget.DateRange = req.DateRange
	widget.DateFrom = req.DateFrom
	widget.DateTo = req.DateTo
	widget.SplitBy = req.SplitBy
	widget.Filters = req.Filters
	if widget.Filters == nil {
		widget.Filters = []Filter{}
	}

	return widget, nil
}

// DeleteWidget deletes a widget.
func (s *Service) DeleteWidget(ctx context.Context, orgID, dashboardID, widgetID uuid.UUID) error {
	// Verify dashboard ownership
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

	// Verify widget exists and belongs to dashboard
	widget, err := s.repo.GetWidgetByID(ctx, widgetID)
	if err != nil {
		return fmt.Errorf("failed to get widget: %w", err)
	}
	if widget == nil {
		return ErrWidgetNotFound
	}
	if widget.DashboardID != dashboardID {
		return ErrWidgetNotFound
	}

	if err := s.repo.DeleteWidget(ctx, widgetID); err != nil {
		return fmt.Errorf("failed to delete widget: %w", err)
	}

	return nil
}

// ReorderWidgets reorders widgets on a dashboard.
func (s *Service) ReorderWidgets(ctx context.Context, orgID, dashboardID uuid.UUID, widgetIDs []uuid.UUID) error {
	// Verify dashboard ownership
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

	if err := s.repo.UpdateWidgetPositions(ctx, dashboardID, widgetIDs); err != nil {
		return fmt.Errorf("failed to reorder widgets: %w", err)
	}

	return nil
}
