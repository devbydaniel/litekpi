package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// DB wraps a pgx connection pool.
type DB struct {
	Pool *pgxpool.Pool
}

// New creates a new database connection pool.
func New(ctx context.Context, databaseURL string) (*DB, error) {
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("parsing database URL: %w", err)
	}

	// Configure connection pool
	config.MaxConns = 25
	config.MinConns = 5
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = 30 * time.Minute
	config.HealthCheckPeriod = time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("creating connection pool: %w", err)
	}

	// Verify connection
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("pinging database: %w", err)
	}

	return &DB{Pool: pool}, nil
}

// Close closes the database connection pool.
func (db *DB) Close() {
	if db.Pool != nil {
		db.Pool.Close()
	}
}

// Health checks if the database connection is healthy.
func (db *DB) Health(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := db.Pool.Ping(ctx); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

	return nil
}
