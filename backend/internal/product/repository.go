package product

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository handles database operations for products.
type Repository struct {
	pool *pgxpool.Pool
}

// NewRepository creates a new product repository.
func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

// CreateProduct creates a new product.
func (r *Repository) CreateProduct(ctx context.Context, orgID uuid.UUID, name, apiKeyHash string) (*Product, error) {
	product := &Product{
		ID:             uuid.New(),
		Name:           name,
		OrganizationID: orgID,
		APIKeyHash:     apiKeyHash,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	_, err := r.pool.Exec(ctx,
		`INSERT INTO products (id, name, organization_id, api_key_hash, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		product.ID, product.Name, product.OrganizationID, product.APIKeyHash, product.CreatedAt, product.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return product, nil
}

// GetProductByID retrieves a product by its ID.
func (r *Repository) GetProductByID(ctx context.Context, id uuid.UUID) (*Product, error) {
	product := &Product{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, name, organization_id, api_key_hash, created_at, updated_at
		FROM products WHERE id = $1`,
		id,
	).Scan(&product.ID, &product.Name, &product.OrganizationID, &product.APIKeyHash, &product.CreatedAt, &product.UpdatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return product, nil
}

// GetProductsByOrganizationID retrieves all products for an organization.
func (r *Repository) GetProductsByOrganizationID(ctx context.Context, orgID uuid.UUID) ([]Product, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, name, organization_id, api_key_hash, created_at, updated_at
		FROM products WHERE organization_id = $1
		ORDER BY created_at DESC`,
		orgID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.Name, &p.OrganizationID, &p.APIKeyHash, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

// UpdateAPIKeyHash updates the API key hash for a product.
func (r *Repository) UpdateAPIKeyHash(ctx context.Context, id uuid.UUID, newHash string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE products SET api_key_hash = $1 WHERE id = $2`,
		newHash, id,
	)
	return err
}

// DeleteProduct deletes a product by its ID.
func (r *Repository) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx,
		`DELETE FROM products WHERE id = $1`,
		id,
	)
	return err
}

// GetProductByAPIKeyHash retrieves a product by its API key hash.
func (r *Repository) GetProductByAPIKeyHash(ctx context.Context, keyHash string) (*Product, error) {
	product := &Product{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, name, organization_id, api_key_hash, created_at, updated_at
		FROM products WHERE api_key_hash = $1`,
		keyHash,
	).Scan(&product.ID, &product.Name, &product.OrganizationID, &product.APIKeyHash, &product.CreatedAt, &product.UpdatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return product, nil
}
