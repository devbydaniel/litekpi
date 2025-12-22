package timeseries

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository handles database operations for time series.
type Repository struct {
	pool *pgxpool.Pool
}

// NewRepository creates a new time series repository.
func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

// Create creates a new time series.
func (r *Repository) Create(ctx context.Context, dashboardID, dataSourceID uuid.UUID, req CreateTimeSeriesRequest, position int) (*TimeSeries, error) {
	filtersJSON, err := json.Marshal(req.Filters)
	if err != nil {
		return nil, err
	}
	if req.Filters == nil {
		filtersJSON = []byte("[]")
	}

	ts := &TimeSeries{
		ID:              uuid.New(),
		DashboardID:     dashboardID,
		DataSourceID:    dataSourceID,
		Title:           req.Title,
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

	if ts.Filters == nil {
		ts.Filters = []Filter{}
	}

	_, err = r.pool.Exec(ctx,
		`INSERT INTO time_series (id, dashboard_id, data_source_id, title, measurement_name, chart_type, date_range, date_from, date_to, split_by, filters, position, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`,
		ts.ID, ts.DashboardID, ts.DataSourceID, ts.Title, ts.MeasurementName, ts.ChartType, ts.DateRange, ts.DateFrom, ts.DateTo, ts.SplitBy, filtersJSON, ts.Position, ts.CreatedAt, ts.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return ts, nil
}

// GetByID retrieves a time series by its ID.
func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*TimeSeries, error) {
	ts := &TimeSeries{}
	var filtersJSON []byte
	var chartType string
	err := r.pool.QueryRow(ctx,
		`SELECT id, dashboard_id, data_source_id, title, measurement_name, chart_type, date_range, date_from, date_to, split_by, filters, position, created_at, updated_at
		FROM time_series WHERE id = $1`,
		id,
	).Scan(&ts.ID, &ts.DashboardID, &ts.DataSourceID, &ts.Title, &ts.MeasurementName, &chartType, &ts.DateRange, &ts.DateFrom, &ts.DateTo, &ts.SplitBy, &filtersJSON, &ts.Position, &ts.CreatedAt, &ts.UpdatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	ts.ChartType = ChartType(chartType)

	if err := json.Unmarshal(filtersJSON, &ts.Filters); err != nil {
		ts.Filters = []Filter{}
	}

	return ts, nil
}

// GetByDashboardID retrieves all time series for a dashboard.
func (r *Repository) GetByDashboardID(ctx context.Context, dashboardID uuid.UUID) ([]TimeSeries, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, dashboard_id, data_source_id, title, measurement_name, chart_type, date_range, date_from, date_to, split_by, filters, position, created_at, updated_at
		FROM time_series WHERE dashboard_id = $1
		ORDER BY position ASC`,
		dashboardID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var timeSeries []TimeSeries
	for rows.Next() {
		var ts TimeSeries
		var filtersJSON []byte
		var chartType string
		if err := rows.Scan(&ts.ID, &ts.DashboardID, &ts.DataSourceID, &ts.Title, &ts.MeasurementName, &chartType, &ts.DateRange, &ts.DateFrom, &ts.DateTo, &ts.SplitBy, &filtersJSON, &ts.Position, &ts.CreatedAt, &ts.UpdatedAt); err != nil {
			return nil, err
		}
		ts.ChartType = ChartType(chartType)
		if err := json.Unmarshal(filtersJSON, &ts.Filters); err != nil {
			ts.Filters = []Filter{}
		}
		timeSeries = append(timeSeries, ts)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if timeSeries == nil {
		timeSeries = []TimeSeries{}
	}

	return timeSeries, nil
}

// Update updates a time series configuration.
func (r *Repository) Update(ctx context.Context, id uuid.UUID, req UpdateTimeSeriesRequest) error {
	filtersJSON, err := json.Marshal(req.Filters)
	if err != nil {
		return err
	}
	if req.Filters == nil {
		filtersJSON = []byte("[]")
	}

	_, err = r.pool.Exec(ctx,
		`UPDATE time_series SET title = $1, chart_type = $2, date_range = $3, date_from = $4, date_to = $5, split_by = $6, filters = $7, updated_at = NOW() WHERE id = $8`,
		req.Title, req.ChartType, req.DateRange, req.DateFrom, req.DateTo, req.SplitBy, filtersJSON, id,
	)
	return err
}

// Delete deletes a time series by its ID.
func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx,
		`DELETE FROM time_series WHERE id = $1`,
		id,
	)
	return err
}

// GetMaxPosition gets the maximum position for time series in a dashboard.
func (r *Repository) GetMaxPosition(ctx context.Context, dashboardID uuid.UUID) (int, error) {
	var maxPos *int
	err := r.pool.QueryRow(ctx,
		`SELECT MAX(position) FROM time_series WHERE dashboard_id = $1`,
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

// UpdatePositions updates the positions of multiple time series.
func (r *Repository) UpdatePositions(ctx context.Context, dashboardID uuid.UUID, timeSeriesIDs []uuid.UUID) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	for i, tsID := range timeSeriesIDs {
		_, err := tx.Exec(ctx,
			`UPDATE time_series SET position = $1, updated_at = NOW() WHERE id = $2 AND dashboard_id = $3`,
			i, tsID, dashboardID,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}
