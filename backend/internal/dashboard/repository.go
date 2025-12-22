package dashboard

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository handles database operations for dashboards.
type Repository struct {
	pool *pgxpool.Pool
}

// NewRepository creates a new dashboard repository.
func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

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
