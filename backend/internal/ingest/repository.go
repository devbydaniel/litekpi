package ingest

import (
	"context"
	"encoding/json"
	"errors"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository handles database operations for measurements.
type Repository struct {
	pool *pgxpool.Pool
}

// NewRepository creates a new ingest repository.
func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

// CreateMeasurement creates a single measurement in the database.
func (r *Repository) CreateMeasurement(ctx context.Context, productID uuid.UUID, name string, value float64, timestamp time.Time, metadata map[string]string) (*Measurement, error) {
	measurement := &Measurement{
		ID:        uuid.New(),
		ProductID: productID,
		Name:      name,
		Value:     value,
		Timestamp: timestamp,
		Metadata:  metadata,
		CreatedAt: time.Now(),
	}

	_, err := r.pool.Exec(ctx,
		`INSERT INTO measurements (id, product_id, name, value, timestamp, metadata, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		measurement.ID, measurement.ProductID, measurement.Name, measurement.Value, measurement.Timestamp, measurement.Metadata, measurement.CreatedAt,
	)
	if err != nil {
		// Check for unique constraint violation (duplicate measurement)
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, ErrDuplicateMeasurement
		}
		return nil, err
	}

	return measurement, nil
}

// CreateMeasurementsBatch creates multiple measurements in a single transaction.
// Returns the count of inserted measurements or an error.
func (r *Repository) CreateMeasurementsBatch(ctx context.Context, productID uuid.UUID, requests []IngestRequest, timestamps []time.Time) (int, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	count := 0
	for i, req := range requests {
		_, err := tx.Exec(ctx,
			`INSERT INTO measurements (id, product_id, name, value, timestamp, metadata, created_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7)`,
			uuid.New(), productID, req.Name, req.Value, timestamps[i], req.Metadata, time.Now(),
		)
		if err != nil {
			// Check for unique constraint violation (duplicate measurement)
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23505" {
				return 0, ErrDuplicateMeasurement
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

// GetMeasurementByID retrieves a measurement by its ID.
func (r *Repository) GetMeasurementByID(ctx context.Context, id uuid.UUID) (*Measurement, error) {
	measurement := &Measurement{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, product_id, name, value, timestamp, metadata, created_at
		FROM measurements WHERE id = $1`,
		id,
	).Scan(&measurement.ID, &measurement.ProductID, &measurement.Name, &measurement.Value, &measurement.Timestamp, &measurement.Metadata, &measurement.CreatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return measurement, nil
}

// GetMeasurementNames retrieves distinct measurement names with their metadata keys for a product.
func (r *Repository) GetMeasurementNames(ctx context.Context, productID uuid.UUID) ([]MeasurementSummary, error) {
	// Query to get distinct names and all unique metadata keys per name
	rows, err := r.pool.Query(ctx,
		`SELECT 
			name,
			COALESCE(
				array_agg(DISTINCT key ORDER BY key) FILTER (WHERE key IS NOT NULL),
				'{}'::text[]
			) as metadata_keys
		FROM measurements
		LEFT JOIN LATERAL (
			SELECT jsonb_object_keys(metadata) as key
			WHERE metadata IS NOT NULL AND metadata != 'null'::jsonb
		) keys ON true
		WHERE product_id = $1
		GROUP BY name
		ORDER BY name`,
		productID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var summaries []MeasurementSummary
	for rows.Next() {
		var summary MeasurementSummary
		if err := rows.Scan(&summary.Name, &summary.MetadataKeys); err != nil {
			return nil, err
		}
		summaries = append(summaries, summary)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if summaries == nil {
		summaries = []MeasurementSummary{}
	}

	return summaries, nil
}

// GetMetadataValues retrieves all unique metadata key-value combinations for a specific measurement.
func (r *Repository) GetMetadataValues(ctx context.Context, productID uuid.UUID, measurementName string) ([]MetadataValues, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT DISTINCT metadata
		FROM measurements
		WHERE product_id = $1 AND name = $2 AND metadata IS NOT NULL AND metadata != 'null'::jsonb`,
		productID, measurementName,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Collect all metadata values
	keyValues := make(map[string]map[string]struct{})
	for rows.Next() {
		var metadataJSON []byte
		if err := rows.Scan(&metadataJSON); err != nil {
			return nil, err
		}

		var metadata map[string]string
		if err := json.Unmarshal(metadataJSON, &metadata); err != nil {
			continue // Skip invalid metadata
		}

		for k, v := range metadata {
			if keyValues[k] == nil {
				keyValues[k] = make(map[string]struct{})
			}
			keyValues[k][v] = struct{}{}
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Convert to response format
	var result []MetadataValues
	for key, valuesSet := range keyValues {
		values := make([]string, 0, len(valuesSet))
		for v := range valuesSet {
			values = append(values, v)
		}
		sort.Strings(values)
		result = append(result, MetadataValues{
			Key:    key,
			Values: values,
		})
	}

	// Sort by key
	sort.Slice(result, func(i, j int) bool {
		return result[i].Key < result[j].Key
	})

	if result == nil {
		result = []MetadataValues{}
	}

	return result, nil
}

// GetAggregatedMeasurements retrieves daily aggregated values with optional metadata filtering.
func (r *Repository) GetAggregatedMeasurements(ctx context.Context, productID uuid.UUID, name string, startDate, endDate time.Time, metadataFilters map[string]string) ([]AggregatedDataPoint, error) {
	// Build the query with optional metadata filtering
	query := `SELECT 
		DATE(timestamp) as date,
		SUM(value) as sum,
		COUNT(*) as count
	FROM measurements
	WHERE product_id = $1 AND name = $2 AND timestamp >= $3 AND timestamp < $4`

	args := []interface{}{productID, name, startDate, endDate}

	// Add metadata filters using JSONB containment operator
	if len(metadataFilters) > 0 {
		filterJSON, err := json.Marshal(metadataFilters)
		if err != nil {
			return nil, err
		}
		query += ` AND metadata @> $5`
		args = append(args, filterJSON)
	}

	query += ` GROUP BY DATE(timestamp) ORDER BY DATE(timestamp)`

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
		dp.Date = date.Format("2006-01-02")
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

// GetAggregatedMeasurementsSplitBy retrieves daily aggregated values split by a metadata key.
// Returns raw data without any top-N aggregation (that's handled by the service layer).
func (r *Repository) GetAggregatedMeasurementsSplitBy(ctx context.Context, productID uuid.UUID, name string, startDate, endDate time.Time, metadataFilters map[string]string, splitByKey string) ([]SplitSeries, error) {
	// Build the query with split by metadata key
	query := `SELECT
		metadata->>$5 as split_key,
		DATE(timestamp) as date,
		SUM(value) as sum,
		COUNT(*) as count
	FROM measurements
	WHERE product_id = $1 AND name = $2 AND timestamp >= $3 AND timestamp < $4
	  AND metadata ? $5`

	args := []interface{}{productID, name, startDate, endDate, splitByKey}

	// Add additional metadata filters using JSONB containment operator
	if len(metadataFilters) > 0 {
		filterJSON, err := json.Marshal(metadataFilters)
		if err != nil {
			return nil, err
		}
		query += ` AND metadata @> $6`
		args = append(args, filterJSON)
	}

	query += ` GROUP BY split_key, DATE(timestamp) ORDER BY split_key, date`

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Collect data points grouped by split key
	seriesMap := make(map[string][]AggregatedDataPoint)
	for rows.Next() {
		var splitKey string
		var dp AggregatedDataPoint
		var date time.Time
		if err := rows.Scan(&splitKey, &date, &dp.Sum, &dp.Count); err != nil {
			return nil, err
		}
		dp.Date = date.Format("2006-01-02")
		seriesMap[splitKey] = append(seriesMap[splitKey], dp)
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

// GetPreferences retrieves saved chart preferences for a measurement.
// Returns nil if no preferences are saved.
func (r *Repository) GetPreferences(ctx context.Context, productID uuid.UUID, measurementName string) (*MeasurementPreferences, error) {
	var preferencesJSON []byte
	err := r.pool.QueryRow(ctx,
		`SELECT preferences FROM measurement_preferences
		WHERE product_id = $1 AND measurement_name = $2`,
		productID, measurementName,
	).Scan(&preferencesJSON)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var prefs MeasurementPreferences
	if err := json.Unmarshal(preferencesJSON, &prefs); err != nil {
		return nil, err
	}

	return &prefs, nil
}

// SavePreferences saves or updates chart preferences for a measurement.
func (r *Repository) SavePreferences(ctx context.Context, productID uuid.UUID, measurementName string, prefs *MeasurementPreferences) error {
	preferencesJSON, err := json.Marshal(prefs)
	if err != nil {
		return err
	}

	_, err = r.pool.Exec(ctx,
		`INSERT INTO measurement_preferences (id, product_id, measurement_name, preferences, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		ON CONFLICT (product_id, measurement_name)
		DO UPDATE SET preferences = $4, updated_at = NOW()`,
		uuid.New(), productID, measurementName, preferencesJSON,
	)
	return err
}
