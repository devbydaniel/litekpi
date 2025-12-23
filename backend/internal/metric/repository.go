package metric

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository handles database operations for metrics.
type Repository struct {
	pool *pgxpool.Pool
}

// NewRepository creates a new metric repository.
func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

// Create creates a new metric.
func (r *Repository) Create(ctx context.Context, dashboardID, dataSourceID uuid.UUID, req CreateMetricRequest, position int) (*Metric, error) {
	filtersJSON, err := json.Marshal(req.Filters)
	if err != nil {
		return nil, err
	}
	if req.Filters == nil {
		filtersJSON = []byte("[]")
	}

	m := &Metric{
		ID:                    uuid.New(),
		DashboardID:           dashboardID,
		DataSourceID:          dataSourceID,
		Label:                 req.Label,
		MeasurementName:       req.MeasurementName,
		Timeframe:             req.Timeframe,
		DateFrom:              req.DateFrom,
		DateTo:                req.DateTo,
		Filters:               req.Filters,
		Aggregation:           req.Aggregation,
		AggregationKey:        req.AggregationKey,
		Granularity:           req.Granularity,
		DisplayMode:           req.DisplayMode,
		ComparisonEnabled:     req.ComparisonEnabled,
		ComparisonDisplayType: req.ComparisonDisplayType,
		ChartType:             req.ChartType,
		SplitBy:               req.SplitBy,
		Position:              position,
		CreatedAt:             time.Now(),
		UpdatedAt:             time.Now(),
	}

	if m.Filters == nil {
		m.Filters = []Filter{}
	}

	_, err = r.pool.Exec(ctx,
		`INSERT INTO metrics (id, dashboard_id, data_source_id, label, measurement_name, timeframe, date_from, date_to, filters, aggregation, aggregation_key, granularity, display_mode, comparison_enabled, comparison_display_type, chart_type, split_by, position, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20)`,
		m.ID, m.DashboardID, m.DataSourceID, m.Label, m.MeasurementName, m.Timeframe, m.DateFrom, m.DateTo, filtersJSON, m.Aggregation, m.AggregationKey, m.Granularity, m.DisplayMode, m.ComparisonEnabled, m.ComparisonDisplayType, m.ChartType, m.SplitBy, m.Position, m.CreatedAt, m.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return m, nil
}

// GetByID retrieves a metric by its ID.
func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*Metric, error) {
	m := &Metric{}
	var filtersJSON []byte
	var aggregation, displayMode string
	var granularity *string
	var comparisonDisplayType, chartType *string
	err := r.pool.QueryRow(ctx,
		`SELECT id, dashboard_id, data_source_id, label, measurement_name, timeframe, date_from, date_to, filters, aggregation, aggregation_key, granularity, display_mode, comparison_enabled, comparison_display_type, chart_type, split_by, position, created_at, updated_at
		FROM metrics WHERE id = $1`,
		id,
	).Scan(&m.ID, &m.DashboardID, &m.DataSourceID, &m.Label, &m.MeasurementName, &m.Timeframe, &m.DateFrom, &m.DateTo, &filtersJSON, &aggregation, &m.AggregationKey, &granularity, &displayMode, &m.ComparisonEnabled, &comparisonDisplayType, &chartType, &m.SplitBy, &m.Position, &m.CreatedAt, &m.UpdatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	m.Aggregation = Aggregation(aggregation)
	if granularity != nil {
		g := Granularity(*granularity)
		m.Granularity = &g
	}
	m.DisplayMode = DisplayMode(displayMode)
	if comparisonDisplayType != nil {
		cdt := ComparisonDisplayType(*comparisonDisplayType)
		m.ComparisonDisplayType = &cdt
	}
	if chartType != nil {
		ct := ChartType(*chartType)
		m.ChartType = &ct
	}

	if err := json.Unmarshal(filtersJSON, &m.Filters); err != nil {
		m.Filters = []Filter{}
	}

	return m, nil
}

// GetByDashboardID retrieves all metrics for a dashboard.
func (r *Repository) GetByDashboardID(ctx context.Context, dashboardID uuid.UUID) ([]Metric, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, dashboard_id, data_source_id, label, measurement_name, timeframe, date_from, date_to, filters, aggregation, aggregation_key, granularity, display_mode, comparison_enabled, comparison_display_type, chart_type, split_by, position, created_at, updated_at
		FROM metrics WHERE dashboard_id = $1
		ORDER BY position ASC`,
		dashboardID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metrics []Metric
	for rows.Next() {
		var m Metric
		var filtersJSON []byte
		var aggregation, displayMode string
		var granularity *string
		var comparisonDisplayType, chartType *string
		if err := rows.Scan(&m.ID, &m.DashboardID, &m.DataSourceID, &m.Label, &m.MeasurementName, &m.Timeframe, &m.DateFrom, &m.DateTo, &filtersJSON, &aggregation, &m.AggregationKey, &granularity, &displayMode, &m.ComparisonEnabled, &comparisonDisplayType, &chartType, &m.SplitBy, &m.Position, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, err
		}
		m.Aggregation = Aggregation(aggregation)
		if granularity != nil {
			g := Granularity(*granularity)
			m.Granularity = &g
		}
		m.DisplayMode = DisplayMode(displayMode)
		if comparisonDisplayType != nil {
			cdt := ComparisonDisplayType(*comparisonDisplayType)
			m.ComparisonDisplayType = &cdt
		}
		if chartType != nil {
			ct := ChartType(*chartType)
			m.ChartType = &ct
		}
		if err := json.Unmarshal(filtersJSON, &m.Filters); err != nil {
			m.Filters = []Filter{}
		}
		metrics = append(metrics, m)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if metrics == nil {
		metrics = []Metric{}
	}

	return metrics, nil
}

// Update updates a metric's configuration.
func (r *Repository) Update(ctx context.Context, id uuid.UUID, req UpdateMetricRequest) error {
	filtersJSON, err := json.Marshal(req.Filters)
	if err != nil {
		return err
	}
	if req.Filters == nil {
		filtersJSON = []byte("[]")
	}

	_, err = r.pool.Exec(ctx,
		`UPDATE metrics SET label = $1, timeframe = $2, date_from = $3, date_to = $4, filters = $5, aggregation = $6, aggregation_key = $7, granularity = $8, display_mode = $9, comparison_enabled = $10, comparison_display_type = $11, chart_type = $12, split_by = $13, updated_at = NOW() WHERE id = $14`,
		req.Label, req.Timeframe, req.DateFrom, req.DateTo, filtersJSON, req.Aggregation, req.AggregationKey, req.Granularity, req.DisplayMode, req.ComparisonEnabled, req.ComparisonDisplayType, req.ChartType, req.SplitBy, id,
	)
	return err
}

// Delete deletes a metric by its ID.
func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx,
		`DELETE FROM metrics WHERE id = $1`,
		id,
	)
	return err
}

// GetMaxPosition gets the maximum position for metrics in a dashboard.
func (r *Repository) GetMaxPosition(ctx context.Context, dashboardID uuid.UUID) (int, error) {
	var maxPos *int
	err := r.pool.QueryRow(ctx,
		`SELECT MAX(position) FROM metrics WHERE dashboard_id = $1`,
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

// UpdatePositions updates the positions of multiple metrics.
func (r *Repository) UpdatePositions(ctx context.Context, dashboardID uuid.UUID, metricIDs []uuid.UUID) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	for i, mID := range metricIDs {
		_, err := tx.Exec(ctx,
			`UPDATE metrics SET position = $1, updated_at = NOW() WHERE id = $2 AND dashboard_id = $3`,
			i, mID, dashboardID,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

// Aggregation queries - these query the measurements table directly

// GetAggregatedMeasurements retrieves aggregated values with optional metadata filtering and granularity.
func (r *Repository) GetAggregatedMeasurements(ctx context.Context, dataSourceID uuid.UUID, name string, startDate, endDate time.Time, metadataFilters map[string]string, granularity Granularity) ([]AggregatedDataPoint, error) {
	dateTrunc := granularityToDateTrunc(granularity)

	query := fmt.Sprintf(`SELECT
		%s as date,
		SUM(value) as sum,
		COUNT(*) as count
	FROM measurements
	WHERE data_source_id = $1 AND name = $2 AND timestamp >= $3 AND timestamp < $4`, dateTrunc)

	args := []interface{}{dataSourceID, name, startDate, endDate}

	// Add metadata filters using JSONB containment operator
	if len(metadataFilters) > 0 {
		filterJSON, err := json.Marshal(metadataFilters)
		if err != nil {
			return nil, err
		}
		query += ` AND metadata @> $5`
		args = append(args, filterJSON)
	}

	query += fmt.Sprintf(` GROUP BY %s ORDER BY %s`, dateTrunc, dateTrunc)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dataPoints []AggregatedDataPoint
	for rows.Next() {
		var dp AggregatedDataPoint
		var date time.Time
		if err := rows.Scan(&date, &dp.Sum, &dp.Count); err != nil {
			return nil, err
		}
		dp.Date = formatDateByGranularity(date, granularity)
		dataPoints = append(dataPoints, dp)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if dataPoints == nil {
		dataPoints = []AggregatedDataPoint{}
	}

	return dataPoints, nil
}

// GetAggregatedMeasurementsSplitBy retrieves aggregated values split by a metadata key.
func (r *Repository) GetAggregatedMeasurementsSplitBy(ctx context.Context, dataSourceID uuid.UUID, name string, startDate, endDate time.Time, metadataFilters map[string]string, splitByKey string, granularity Granularity) ([]SplitSeries, error) {
	dateTrunc := granularityToDateTrunc(granularity)

	query := fmt.Sprintf(`SELECT
		metadata->>$5 as split_key,
		%s as date,
		SUM(value) as sum,
		COUNT(*) as count
	FROM measurements
	WHERE data_source_id = $1 AND name = $2 AND timestamp >= $3 AND timestamp < $4
	  AND metadata ? $5`, dateTrunc)

	args := []interface{}{dataSourceID, name, startDate, endDate, splitByKey}

	// Add additional metadata filters using JSONB containment operator
	if len(metadataFilters) > 0 {
		filterJSON, err := json.Marshal(metadataFilters)
		if err != nil {
			return nil, err
		}
		query += ` AND metadata @> $6`
		args = append(args, filterJSON)
	}

	query += fmt.Sprintf(` GROUP BY split_key, %s ORDER BY split_key, date`, dateTrunc)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Collect data points grouped by split key
	seriesMap := make(map[string][]DataPoint)
	for rows.Next() {
		var splitKey string
		var dp AggregatedDataPoint
		var date time.Time
		if err := rows.Scan(&splitKey, &date, &dp.Sum, &dp.Count); err != nil {
			return nil, err
		}
		seriesMap[splitKey] = append(seriesMap[splitKey], DataPoint{
			Date:  formatDateByGranularity(date, granularity),
			Value: dp.Sum,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Convert to slice of SplitSeries
	var series []SplitSeries
	for key, dataPoints := range seriesMap {
		series = append(series, SplitSeries{
			Key:        key,
			DataPoints: dataPoints,
		})
	}

	// Sort by key for consistent ordering
	sort.Slice(series, func(i, j int) bool {
		return series[i].Key < series[j].Key
	})

	if series == nil {
		series = []SplitSeries{}
	}

	return series, nil
}

// GetCountUniqueMeasurements retrieves the count of unique values for a metadata key.
func (r *Repository) GetCountUniqueMeasurements(ctx context.Context, dataSourceID uuid.UUID, name string, startDate, endDate time.Time, metadataFilters map[string]string, aggregationKey string, granularity Granularity) ([]AggregatedDataPoint, error) {
	dateTrunc := granularityToDateTrunc(granularity)

	query := fmt.Sprintf(`SELECT
		%s as date,
		COUNT(DISTINCT metadata->>$5) as count
	FROM measurements
	WHERE data_source_id = $1 AND name = $2 AND timestamp >= $3 AND timestamp < $4
	  AND metadata ? $5`, dateTrunc)

	args := []interface{}{dataSourceID, name, startDate, endDate, aggregationKey}

	// Add additional metadata filters using JSONB containment operator
	if len(metadataFilters) > 0 {
		filterJSON, err := json.Marshal(metadataFilters)
		if err != nil {
			return nil, err
		}
		query += ` AND metadata @> $6`
		args = append(args, filterJSON)
	}

	query += fmt.Sprintf(` GROUP BY %s ORDER BY %s`, dateTrunc, dateTrunc)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dataPoints []AggregatedDataPoint
	for rows.Next() {
		var dp AggregatedDataPoint
		var date time.Time
		if err := rows.Scan(&date, &dp.Count); err != nil {
			return nil, err
		}
		dp.Date = formatDateByGranularity(date, granularity)
		dp.Sum = float64(dp.Count) // For count_unique, Sum holds the unique count value
		dataPoints = append(dataPoints, dp)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if dataPoints == nil {
		dataPoints = []AggregatedDataPoint{}
	}

	return dataPoints, nil
}

// Helper functions

func granularityToDateTrunc(g Granularity) string {
	switch g {
	case GranularityWeekly:
		return "DATE_TRUNC('week', timestamp)::date"
	case GranularityMonthly:
		return "DATE_TRUNC('month', timestamp)::date"
	default: // daily
		return "DATE(timestamp)"
	}
}

func formatDateByGranularity(t time.Time, g Granularity) string {
	switch g {
	case GranularityWeekly:
		return t.Format("2006-01-02") // Start of week
	case GranularityMonthly:
		return t.Format("2006-01") // Year-Month
	default: // daily
		return t.Format("2006-01-02")
	}
}

// Scalar aggregation queries - no granularity/grouping, returns single aggregate

// GetScalarAggregate returns the sum and count for the entire timeframe without grouping.
func (r *Repository) GetScalarAggregate(ctx context.Context, dataSourceID uuid.UUID, name string, startDate, endDate time.Time, metadataFilters map[string]string) (sum float64, count int, err error) {
	query := `SELECT COALESCE(SUM(value), 0), COUNT(*)
	FROM measurements
	WHERE data_source_id = $1 AND name = $2 AND timestamp >= $3 AND timestamp < $4`

	args := []interface{}{dataSourceID, name, startDate, endDate}

	if len(metadataFilters) > 0 {
		filterJSON, err := json.Marshal(metadataFilters)
		if err != nil {
			return 0, 0, err
		}
		query += ` AND metadata @> $5`
		args = append(args, filterJSON)
	}

	err = r.pool.QueryRow(ctx, query, args...).Scan(&sum, &count)
	if err != nil {
		return 0, 0, err
	}

	return sum, count, nil
}

// GetScalarCountUnique returns the unique count for the entire timeframe without grouping.
func (r *Repository) GetScalarCountUnique(ctx context.Context, dataSourceID uuid.UUID, name string, startDate, endDate time.Time, metadataFilters map[string]string, aggregationKey string) (int, error) {
	query := `SELECT COUNT(DISTINCT metadata->>$5)
	FROM measurements
	WHERE data_source_id = $1 AND name = $2 AND timestamp >= $3 AND timestamp < $4
	  AND metadata ? $5`

	args := []interface{}{dataSourceID, name, startDate, endDate, aggregationKey}

	if len(metadataFilters) > 0 {
		filterJSON, err := json.Marshal(metadataFilters)
		if err != nil {
			return 0, err
		}
		query += ` AND metadata @> $6`
		args = append(args, filterJSON)
	}

	var count int
	err := r.pool.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
