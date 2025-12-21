package datasource

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository handles database operations for data sources.
type Repository struct {
	pool *pgxpool.Pool
}

// NewRepository creates a new data source repository.
func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

// CreateDataSource creates a new data source.
func (r *Repository) CreateDataSource(ctx context.Context, orgID uuid.UUID, name, apiKeyHash string) (*DataSource, error) {
	ds := &DataSource{
		ID:             uuid.New(),
		Name:           name,
		OrganizationID: orgID,
		APIKeyHash:     apiKeyHash,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	_, err := r.pool.Exec(ctx,
		`INSERT INTO data_sources (id, name, organization_id, api_key_hash, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		ds.ID, ds.Name, ds.OrganizationID, ds.APIKeyHash, ds.CreatedAt, ds.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return ds, nil
}

// GetDataSourceByID retrieves a data source by its ID.
func (r *Repository) GetDataSourceByID(ctx context.Context, id uuid.UUID) (*DataSource, error) {
	ds := &DataSource{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, name, organization_id, api_key_hash, created_at, updated_at
		FROM data_sources WHERE id = $1`,
		id,
	).Scan(&ds.ID, &ds.Name, &ds.OrganizationID, &ds.APIKeyHash, &ds.CreatedAt, &ds.UpdatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return ds, nil
}

// GetDataSourcesByOrganizationID retrieves all data sources for an organization.
func (r *Repository) GetDataSourcesByOrganizationID(ctx context.Context, orgID uuid.UUID) ([]DataSource, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, name, organization_id, api_key_hash, created_at, updated_at
		FROM data_sources WHERE organization_id = $1
		ORDER BY created_at DESC`,
		orgID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dataSources []DataSource
	for rows.Next() {
		var ds DataSource
		if err := rows.Scan(&ds.ID, &ds.Name, &ds.OrganizationID, &ds.APIKeyHash, &ds.CreatedAt, &ds.UpdatedAt); err != nil {
			return nil, err
		}
		dataSources = append(dataSources, ds)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return dataSources, nil
}

// UpdateAPIKeyHash updates the API key hash for a data source.
func (r *Repository) UpdateAPIKeyHash(ctx context.Context, id uuid.UUID, newHash string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE data_sources SET api_key_hash = $1 WHERE id = $2`,
		newHash, id,
	)
	return err
}

// DeleteDataSource deletes a data source by its ID.
func (r *Repository) DeleteDataSource(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx,
		`DELETE FROM data_sources WHERE id = $1`,
		id,
	)
	return err
}

// GetDataSourceByAPIKeyHash retrieves a data source by its API key hash.
func (r *Repository) GetDataSourceByAPIKeyHash(ctx context.Context, keyHash string) (*DataSource, error) {
	ds := &DataSource{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, name, organization_id, api_key_hash, created_at, updated_at
		FROM data_sources WHERE api_key_hash = $1`,
		keyHash,
	).Scan(&ds.ID, &ds.Name, &ds.OrganizationID, &ds.APIKeyHash, &ds.CreatedAt, &ds.UpdatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return ds, nil
}
