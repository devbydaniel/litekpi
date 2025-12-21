package dashboard

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository handles database operations for dashboards and widgets.
type Repository struct {
	pool *pgxpool.Pool
}

// NewRepository creates a new dashboard repository.
func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

// Dashboard operations

// CreateDashboard creates a new dashboard.
func (r *Repository) CreateDashboard(ctx context.Context, orgID uuid.UUID, name string, isDefault bool) (*Dashboard, error) {
	dashboard := &Dashboard{
		ID:             uuid.New(),
		Name:           name,
		OrganizationID: orgID,
		IsDefault:      isDefault,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	_, err := r.pool.Exec(ctx,
		`INSERT INTO dashboards (id, name, organization_id, is_default, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		dashboard.ID, dashboard.Name, dashboard.OrganizationID, dashboard.IsDefault, dashboard.CreatedAt, dashboard.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return dashboard, nil
}

// GetDashboardByID retrieves a dashboard by its ID.
func (r *Repository) GetDashboardByID(ctx context.Context, id uuid.UUID) (*Dashboard, error) {
	dashboard := &Dashboard{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, name, organization_id, is_default, created_at, updated_at
		FROM dashboards WHERE id = $1`,
		id,
	).Scan(&dashboard.ID, &dashboard.Name, &dashboard.OrganizationID, &dashboard.IsDefault, &dashboard.CreatedAt, &dashboard.UpdatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return dashboard, nil
}

// GetDefaultDashboard retrieves the default dashboard for an organization.
func (r *Repository) GetDefaultDashboard(ctx context.Context, orgID uuid.UUID) (*Dashboard, error) {
	dashboard := &Dashboard{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, name, organization_id, is_default, created_at, updated_at
		FROM dashboards WHERE organization_id = $1 AND is_default = TRUE`,
		orgID,
	).Scan(&dashboard.ID, &dashboard.Name, &dashboard.OrganizationID, &dashboard.IsDefault, &dashboard.CreatedAt, &dashboard.UpdatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return dashboard, nil
}

// GetDashboardsByOrganizationID retrieves all dashboards for an organization.
func (r *Repository) GetDashboardsByOrganizationID(ctx context.Context, orgID uuid.UUID) ([]Dashboard, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, name, organization_id, is_default, created_at, updated_at
		FROM dashboards WHERE organization_id = $1
		ORDER BY is_default DESC, created_at ASC`,
		orgID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dashboards []Dashboard
	for rows.Next() {
		var d Dashboard
		if err := rows.Scan(&d.ID, &d.Name, &d.OrganizationID, &d.IsDefault, &d.CreatedAt, &d.UpdatedAt); err != nil {
			return nil, err
		}
		dashboards = append(dashboards, d)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if dashboards == nil {
		dashboards = []Dashboard{}
	}

	return dashboards, nil
}

// UpdateDashboard updates a dashboard's name.
func (r *Repository) UpdateDashboard(ctx context.Context, id uuid.UUID, name string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE dashboards SET name = $1, updated_at = NOW() WHERE id = $2`,
		name, id,
	)
	return err
}

// DeleteDashboard deletes a dashboard by its ID.
func (r *Repository) DeleteDashboard(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx,
		`DELETE FROM dashboards WHERE id = $1`,
		id,
	)
	return err
}

// Widget operations

// CreateWidget creates a new widget.
func (r *Repository) CreateWidget(ctx context.Context, dashboardID, dataSourceID uuid.UUID, req CreateWidgetRequest, position int) (*Widget, error) {
	filtersJSON, err := json.Marshal(req.Filters)
	if err != nil {
		return nil, err
	}
	if req.Filters == nil {
		filtersJSON = []byte("[]")
	}

	widget := &Widget{
		ID:              uuid.New(),
		DashboardID:     dashboardID,
		DataSourceID:    dataSourceID,
		MeasurementName: req.MeasurementName,
		ChartType:       req.ChartType,
		DateRange:       req.DateRange,
		DateFrom:        req.DateFrom,
		DateTo:          req.DateTo,
		SplitBy:         req.SplitBy,
		Filters:         req.Filters,
		Position:        position,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if widget.Filters == nil {
		widget.Filters = []Filter{}
	}

	_, err = r.pool.Exec(ctx,
		`INSERT INTO widgets (id, dashboard_id, data_source_id, measurement_name, chart_type, date_range, date_from, date_to, split_by, filters, position, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`,
		widget.ID, widget.DashboardID, widget.DataSourceID, widget.MeasurementName, widget.ChartType, widget.DateRange, widget.DateFrom, widget.DateTo, widget.SplitBy, filtersJSON, widget.Position, widget.CreatedAt, widget.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return widget, nil
}

// GetWidgetByID retrieves a widget by its ID.
func (r *Repository) GetWidgetByID(ctx context.Context, id uuid.UUID) (*Widget, error) {
	widget := &Widget{}
	var filtersJSON []byte
	err := r.pool.QueryRow(ctx,
		`SELECT id, dashboard_id, data_source_id, measurement_name, chart_type, date_range, date_from, date_to, split_by, filters, position, created_at, updated_at
		FROM widgets WHERE id = $1`,
		id,
	).Scan(&widget.ID, &widget.DashboardID, &widget.DataSourceID, &widget.MeasurementName, &widget.ChartType, &widget.DateRange, &widget.DateFrom, &widget.DateTo, &widget.SplitBy, &filtersJSON, &widget.Position, &widget.CreatedAt, &widget.UpdatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(filtersJSON, &widget.Filters); err != nil {
		widget.Filters = []Filter{}
	}

	return widget, nil
}

// GetWidgetsByDashboardID retrieves all widgets for a dashboard.
func (r *Repository) GetWidgetsByDashboardID(ctx context.Context, dashboardID uuid.UUID) ([]Widget, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, dashboard_id, data_source_id, measurement_name, chart_type, date_range, date_from, date_to, split_by, filters, position, created_at, updated_at
		FROM widgets WHERE dashboard_id = $1
		ORDER BY position ASC`,
		dashboardID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var widgets []Widget
	for rows.Next() {
		var w Widget
		var filtersJSON []byte
		if err := rows.Scan(&w.ID, &w.DashboardID, &w.DataSourceID, &w.MeasurementName, &w.ChartType, &w.DateRange, &w.DateFrom, &w.DateTo, &w.SplitBy, &filtersJSON, &w.Position, &w.CreatedAt, &w.UpdatedAt); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(filtersJSON, &w.Filters); err != nil {
			w.Filters = []Filter{}
		}
		widgets = append(widgets, w)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if widgets == nil {
		widgets = []Widget{}
	}

	return widgets, nil
}

// UpdateWidget updates a widget's configuration.
func (r *Repository) UpdateWidget(ctx context.Context, id uuid.UUID, req UpdateWidgetRequest) error {
	filtersJSON, err := json.Marshal(req.Filters)
	if err != nil {
		return err
	}
	if req.Filters == nil {
		filtersJSON = []byte("[]")
	}

	_, err = r.pool.Exec(ctx,
		`UPDATE widgets SET chart_type = $1, date_range = $2, date_from = $3, date_to = $4, split_by = $5, filters = $6, updated_at = NOW() WHERE id = $7`,
		req.ChartType, req.DateRange, req.DateFrom, req.DateTo, req.SplitBy, filtersJSON, id,
	)
	return err
}

// DeleteWidget deletes a widget by its ID.
func (r *Repository) DeleteWidget(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx,
		`DELETE FROM widgets WHERE id = $1`,
		id,
	)
	return err
}

// GetMaxWidgetPosition gets the maximum position for widgets in a dashboard.
func (r *Repository) GetMaxWidgetPosition(ctx context.Context, dashboardID uuid.UUID) (int, error) {
	var maxPos *int
	err := r.pool.QueryRow(ctx,
		`SELECT MAX(position) FROM widgets WHERE dashboard_id = $1`,
		dashboardID,
	).Scan(&maxPos)

	if err != nil {
		return 0, err
	}
	if maxPos == nil {
		return 0, nil
	}

	return *maxPos, nil
}

// UpdateWidgetPositions updates the positions of multiple widgets.
func (r *Repository) UpdateWidgetPositions(ctx context.Context, dashboardID uuid.UUID, widgetIDs []uuid.UUID) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	for i, widgetID := range widgetIDs {
		_, err := tx.Exec(ctx,
			`UPDATE widgets SET position = $1, updated_at = NOW() WHERE id = $2 AND dashboard_id = $3`,
			i, widgetID, dashboardID,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}
