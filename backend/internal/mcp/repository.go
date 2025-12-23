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

// Create creates a new MCP API key.
func (r *Repository) Create(ctx context.Context, orgID uuid.UUID, name, apiKeyHash string, createdBy uuid.UUID) (*MCPAPIKey, error) {
	key := &MCPAPIKey{
		ID:             uuid.New(),
		OrganizationID: orgID,
		Name:           name,
		APIKeyHash:     apiKeyHash,
		CreatedBy:      createdBy,
		CreatedAt:      time.Now(),
	}

	_, err := r.pool.Exec(ctx,
		`INSERT INTO mcp_api_keys (id, organization_id, name, api_key_hash, created_by, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		key.ID, key.OrganizationID, key.Name, key.APIKeyHash, key.CreatedBy, key.CreatedAt,
	)
	if err != nil {
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

// UpdateLastUsed updates the last_used_at timestamp for an MCP API key.
func (r *Repository) UpdateLastUsed(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE mcp_api_keys SET last_used_at = $1 WHERE id = $2`,
		time.Now(), id,
	)
	return err
}
