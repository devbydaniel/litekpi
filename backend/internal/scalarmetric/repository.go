package scalarmetric

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository handles database operations for scalar metrics.
type Repository struct {
	pool *pgxpool.Pool
}

// NewRepository creates a new scalar metric repository.
func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

// Create creates a new scalar metric.
func (r *Repository) Create(ctx context.Context, dashboardID, dataSourceID uuid.UUID, req CreateScalarMetricRequest, position int) (*ScalarMetric, error) {
	filtersJSON, err := json.Marshal(req.Filters)
	if err != nil {
		return nil, err
	}
	if req.Filters == nil {
		filtersJSON = []byte("[]")
	}

	sm := &ScalarMetric{
		ID:                    uuid.New(),
		DashboardID:           dashboardID,
		DataSourceID:          dataSourceID,
		Label:                 req.Label,
		MeasurementName:       req.MeasurementName,
		Timeframe:             req.Timeframe,
		Aggregation:           req.Aggregation,
		Filters:               req.Filters,
		ComparisonEnabled:     req.ComparisonEnabled,
		ComparisonDisplayType: req.ComparisonDisplayType,
		Position:              position,
		CreatedAt:             time.Now(),
		UpdatedAt:             time.Now(),
	}

	if sm.Filters == nil {
		sm.Filters = []Filter{}
	}

	_, err = r.pool.Exec(ctx,
		`INSERT INTO scalar_metrics (id, dashboard_id, data_source_id, label, measurement_name, timeframe, aggregation, filters, comparison_enabled, comparison_display_type, position, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`,
		sm.ID, sm.DashboardID, sm.DataSourceID, sm.Label, sm.MeasurementName, sm.Timeframe, sm.Aggregation, filtersJSON, sm.ComparisonEnabled, sm.ComparisonDisplayType, sm.Position, sm.CreatedAt, sm.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return sm, nil
}

// GetByID retrieves a scalar metric by its ID.
func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*ScalarMetric, error) {
	sm := &ScalarMetric{}
	var filtersJSON []byte
	var aggregation string
	var comparisonDisplayType *string
	err := r.pool.QueryRow(ctx,
		`SELECT id, dashboard_id, data_source_id, label, measurement_name, timeframe, aggregation, filters, comparison_enabled, comparison_display_type, position, created_at, updated_at
		FROM scalar_metrics WHERE id = $1`,
		id,
	).Scan(&sm.ID, &sm.DashboardID, &sm.DataSourceID, &sm.Label, &sm.MeasurementName, &sm.Timeframe, &aggregation, &filtersJSON, &sm.ComparisonEnabled, &comparisonDisplayType, &sm.Position, &sm.CreatedAt, &sm.UpdatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	sm.Aggregation = Aggregation(aggregation)
	if comparisonDisplayType != nil {
		cdt := ComparisonDisplayType(*comparisonDisplayType)
		sm.ComparisonDisplayType = &cdt
	}

	if err := json.Unmarshal(filtersJSON, &sm.Filters); err != nil {
		sm.Filters = []Filter{}
	}

	return sm, nil
}

// GetByDashboardID retrieves all scalar metrics for a dashboard.
func (r *Repository) GetByDashboardID(ctx context.Context, dashboardID uuid.UUID) ([]ScalarMetric, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, dashboard_id, data_source_id, label, measurement_name, timeframe, aggregation, filters, comparison_enabled, comparison_display_type, position, created_at, updated_at
		FROM scalar_metrics WHERE dashboard_id = $1
		ORDER BY position ASC`,
		dashboardID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var scalarMetrics []ScalarMetric
	for rows.Next() {
		var sm ScalarMetric
		var filtersJSON []byte
		var aggregation string
		var comparisonDisplayType *string
		if err := rows.Scan(&sm.ID, &sm.DashboardID, &sm.DataSourceID, &sm.Label, &sm.MeasurementName, &sm.Timeframe, &aggregation, &filtersJSON, &sm.ComparisonEnabled, &comparisonDisplayType, &sm.Position, &sm.CreatedAt, &sm.UpdatedAt); err != nil {
			return nil, err
		}
		sm.Aggregation = Aggregation(aggregation)
		if comparisonDisplayType != nil {
			cdt := ComparisonDisplayType(*comparisonDisplayType)
			sm.ComparisonDisplayType = &cdt
		}
		if err := json.Unmarshal(filtersJSON, &sm.Filters); err != nil {
			sm.Filters = []Filter{}
		}
		scalarMetrics = append(scalarMetrics, sm)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if scalarMetrics == nil {
		scalarMetrics = []ScalarMetric{}
	}

	return scalarMetrics, nil
}

// Update updates a scalar metric's configuration.
func (r *Repository) Update(ctx context.Context, id uuid.UUID, req UpdateScalarMetricRequest) error {
	filtersJSON, err := json.Marshal(req.Filters)
	if err != nil {
		return err
	}
	if req.Filters == nil {
		filtersJSON = []byte("[]")
	}

	_, err = r.pool.Exec(ctx,
		`UPDATE scalar_metrics SET label = $1, timeframe = $2, aggregation = $3, filters = $4, comparison_enabled = $5, comparison_display_type = $6, updated_at = NOW() WHERE id = $7`,
		req.Label, req.Timeframe, req.Aggregation, filtersJSON, req.ComparisonEnabled, req.ComparisonDisplayType, id,
	)
	return err
}

// Delete deletes a scalar metric by its ID.
func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx,
		`DELETE FROM scalar_metrics WHERE id = $1`,
		id,
	)
	return err
}

// GetMaxPosition gets the maximum position for scalar metrics in a dashboard.
func (r *Repository) GetMaxPosition(ctx context.Context, dashboardID uuid.UUID) (int, error) {
	var maxPos *int
	err := r.pool.QueryRow(ctx,
		`SELECT MAX(position) FROM scalar_metrics WHERE dashboard_id = $1`,
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

// UpdatePositions updates the positions of multiple scalar metrics.
func (r *Repository) UpdatePositions(ctx context.Context, dashboardID uuid.UUID, scalarMetricIDs []uuid.UUID) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	for i, smID := range scalarMetricIDs {
		_, err := tx.Exec(ctx,
			`UPDATE scalar_metrics SET position = $1, updated_at = NOW() WHERE id = $2 AND dashboard_id = $3`,
			i, smID, dashboardID,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}
