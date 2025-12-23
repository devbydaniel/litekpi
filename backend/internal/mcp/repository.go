package mcp

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository handles database operations for MCP API keys.
type Repository struct {
	pool *pgxpool.Pool
}

// NewRepository creates a new MCP repository.
func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

// Create creates a new MCP API key with associated data sources.
func (r *Repository) Create(ctx context.Context, orgID uuid.UUID, name, apiKeyHash string, createdBy uuid.UUID, dataSourceIDs []uuid.UUID) (*MCPAPIKey, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	key := &MCPAPIKey{
		ID:                   uuid.New(),
		OrganizationID:       orgID,
		Name:                 name,
		APIKeyHash:           apiKeyHash,
		CreatedBy:            createdBy,
		CreatedAt:            time.Now(),
		AllowedDataSourceIDs: dataSourceIDs,
	}

	_, err = tx.Exec(ctx,
		`INSERT INTO mcp_api_keys (id, organization_id, name, api_key_hash, created_by, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		key.ID, key.OrganizationID, key.Name, key.APIKeyHash, key.CreatedBy, key.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	// Insert data source associations
	for _, dsID := range dataSourceIDs {
		_, err = tx.Exec(ctx,
			`INSERT INTO mcp_api_key_data_sources (mcp_api_key_id, data_source_id)
			VALUES ($1, $2)`,
			key.ID, dsID,
		)
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return key, nil
}

// GetByID retrieves an MCP API key by its ID.
func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*MCPAPIKey, error) {
	key := &MCPAPIKey{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, organization_id, name, api_key_hash, created_by, last_used_at, created_at
		FROM mcp_api_keys WHERE id = $1`,
		id,
	).Scan(&key.ID, &key.OrganizationID, &key.Name, &key.APIKeyHash, &key.CreatedBy, &key.LastUsedAt, &key.CreatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Load allowed data source IDs
	key.AllowedDataSourceIDs, err = r.getDataSourceIDsForKey(ctx, key.ID)
	if err != nil {
		return nil, err
	}

	return key, nil
}

// GetByAPIKeyHash retrieves an MCP API key by its hash.
func (r *Repository) GetByAPIKeyHash(ctx context.Context, keyHash string) (*MCPAPIKey, error) {
	key := &MCPAPIKey{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, organization_id, name, api_key_hash, created_by, last_used_at, created_at
		FROM mcp_api_keys WHERE api_key_hash = $1`,
		keyHash,
	).Scan(&key.ID, &key.OrganizationID, &key.Name, &key.APIKeyHash, &key.CreatedBy, &key.LastUsedAt, &key.CreatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Load allowed data source IDs
	key.AllowedDataSourceIDs, err = r.getDataSourceIDsForKey(ctx, key.ID)
	if err != nil {
		return nil, err
	}

	return key, nil
}

// GetByOrganizationID retrieves all MCP API keys for an organization.
func (r *Repository) GetByOrganizationID(ctx context.Context, orgID uuid.UUID) ([]MCPAPIKey, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, organization_id, name, api_key_hash, created_by, last_used_at, created_at
		FROM mcp_api_keys WHERE organization_id = $1
		ORDER BY created_at DESC`,
		orgID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var keys []MCPAPIKey
	for rows.Next() {
		var key MCPAPIKey
		if err := rows.Scan(&key.ID, &key.OrganizationID, &key.Name, &key.APIKeyHash, &key.CreatedBy, &key.LastUsedAt, &key.CreatedAt); err != nil {
			return nil, err
		}
		keys = append(keys, key)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Load allowed data source IDs for each key
	for i := range keys {
		keys[i].AllowedDataSourceIDs, err = r.getDataSourceIDsForKey(ctx, keys[i].ID)
		if err != nil {
			return nil, err
		}
	}

	return keys, nil
}

// Delete deletes an MCP API key by its ID.
func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx,
		`DELETE FROM mcp_api_keys WHERE id = $1`,
		id,
	)
	return err
}

// UpdateDataSources updates the allowed data sources for an MCP API key.
func (r *Repository) UpdateDataSources(ctx context.Context, keyID uuid.UUID, dataSourceIDs []uuid.UUID) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Delete existing associations
	_, err = tx.Exec(ctx,
		`DELETE FROM mcp_api_key_data_sources WHERE mcp_api_key_id = $1`,
		keyID,
	)
	if err != nil {
		return err
	}

	// Insert new associations
	for _, dsID := range dataSourceIDs {
		_, err = tx.Exec(ctx,
			`INSERT INTO mcp_api_key_data_sources (mcp_api_key_id, data_source_id)
			VALUES ($1, $2)`,
			keyID, dsID,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

// UpdateLastUsed updates the last_used_at timestamp for an MCP API key.
func (r *Repository) UpdateLastUsed(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE mcp_api_keys SET last_used_at = $1 WHERE id = $2`,
		time.Now(), id,
	)
	return err
}

// getDataSourceIDsForKey retrieves all allowed data source IDs for an MCP API key.
func (r *Repository) getDataSourceIDsForKey(ctx context.Context, keyID uuid.UUID) ([]uuid.UUID, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT data_source_id FROM mcp_api_key_data_sources WHERE mcp_api_key_id = $1`,
		keyID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	return ids, rows.Err()
}
