package report

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository handles database operations for reports.
type Repository struct {
	pool *pgxpool.Pool
}

// NewRepository creates a new report repository.
func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

// CreateReport creates a new report.
func (r *Repository) CreateReport(ctx context.Context, orgID uuid.UUID, name string) (*Report, error) {
	report := &Report{
		ID:             uuid.New(),
		Name:           name,
		OrganizationID: orgID,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	_, err := r.pool.Exec(ctx,
		`INSERT INTO reports (id, name, organization_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)`,
		report.ID, report.Name, report.OrganizationID, report.CreatedAt, report.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return report, nil
}

// GetReportByID retrieves a report by its ID.
func (r *Repository) GetReportByID(ctx context.Context, id uuid.UUID) (*Report, error) {
	report := &Report{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, name, organization_id, created_at, updated_at
		FROM reports WHERE id = $1`,
		id,
	).Scan(&report.ID, &report.Name, &report.OrganizationID, &report.CreatedAt, &report.UpdatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return report, nil
}

// GetReportsByOrganizationID retrieves all reports for an organization.
func (r *Repository) GetReportsByOrganizationID(ctx context.Context, orgID uuid.UUID) ([]Report, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, name, organization_id, created_at, updated_at
		FROM reports WHERE organization_id = $1
		ORDER BY created_at DESC`,
		orgID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reports []Report
	for rows.Next() {
		var rpt Report
		if err := rows.Scan(&rpt.ID, &rpt.Name, &rpt.OrganizationID, &rpt.CreatedAt, &rpt.UpdatedAt); err != nil {
			return nil, err
		}
		reports = append(reports, rpt)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if reports == nil {
		reports = []Report{}
	}

	return reports, nil
}

// UpdateReport updates a report's name.
func (r *Repository) UpdateReport(ctx context.Context, id uuid.UUID, name string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE reports SET name = $1, updated_at = NOW() WHERE id = $2`,
		name, id,
	)
	return err
}

// DeleteReport deletes a report by its ID.
func (r *Repository) DeleteReport(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx,
		`DELETE FROM reports WHERE id = $1`,
		id,
	)
	return err
}
