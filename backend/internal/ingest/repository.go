package ingest

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository handles database operations for metrics.
type Repository struct {
	pool *pgxpool.Pool
}

// NewRepository creates a new ingest repository.
func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

// CreateMetric creates a single metric in the database.
func (r *Repository) CreateMetric(ctx context.Context, productID uuid.UUID, name string, value float64, timestamp time.Time, metadata map[string]string) (*Metric, error) {
	metric := &Metric{
		ID:        uuid.New(),
		ProductID: productID,
		Name:      name,
		Value:     value,
		Timestamp: timestamp,
		Metadata:  metadata,
		CreatedAt: time.Now(),
	}

	_, err := r.pool.Exec(ctx,
		`INSERT INTO metrics (id, product_id, name, value, timestamp, metadata, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		metric.ID, metric.ProductID, metric.Name, metric.Value, metric.Timestamp, metric.Metadata, metric.CreatedAt,
	)
	if err != nil {
		// Check for unique constraint violation (duplicate metric)
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, ErrDuplicateMetric
		}
		return nil, err
	}

	return metric, nil
}

// CreateMetricsBatch creates multiple metrics in a single transaction.
// Returns the count of inserted metrics or an error.
func (r *Repository) CreateMetricsBatch(ctx context.Context, productID uuid.UUID, metrics []IngestRequest, timestamps []time.Time) (int, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	count := 0
	for i, m := range metrics {
		_, err := tx.Exec(ctx,
			`INSERT INTO metrics (id, product_id, name, value, timestamp, metadata, created_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7)`,
			uuid.New(), productID, m.Name, m.Value, timestamps[i], m.Metadata, time.Now(),
		)
		if err != nil {
			// Check for unique constraint violation (duplicate metric)
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23505" {
				return 0, ErrDuplicateMetric
			}
			return 0, err
		}
		count++
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, err
	}

	return count, nil
}

// GetMetricByID retrieves a metric by its ID.
func (r *Repository) GetMetricByID(ctx context.Context, id uuid.UUID) (*Metric, error) {
	metric := &Metric{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, product_id, name, value, timestamp, metadata, created_at
		FROM metrics WHERE id = $1`,
		id,
	).Scan(&metric.ID, &metric.ProductID, &metric.Name, &metric.Value, &metric.Timestamp, &metric.Metadata, &metric.CreatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return metric, nil
}
