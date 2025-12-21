package kpi

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository handles database operations for KPIs.
type Repository struct {
	pool *pgxpool.Pool
}

// NewRepository creates a new KPI repository.
func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

// CreateKPIForDashboard creates a new KPI for a dashboard.
func (r *Repository) CreateKPIForDashboard(ctx context.Context, dashboardID, dataSourceID uuid.UUID, req CreateKPIRequest, position int) (*KPI, error) {
	return r.createKPI(ctx, &dashboardID, nil, dataSourceID, req, position)
}

// CreateKPIForReport creates a new KPI for a report.
func (r *Repository) CreateKPIForReport(ctx context.Context, reportID, dataSourceID uuid.UUID, req CreateKPIRequest, position int) (*KPI, error) {
	return r.createKPI(ctx, nil, &reportID, dataSourceID, req, position)
}

func (r *Repository) createKPI(ctx context.Context, dashboardID, reportID *uuid.UUID, dataSourceID uuid.UUID, req CreateKPIRequest, position int) (*KPI, error) {
	filtersJSON, err := json.Marshal(req.Filters)
	if err != nil {
		return nil, err
	}
	if req.Filters == nil {
		filtersJSON = []byte("[]")
	}

	kpi := &KPI{
		ID:                    uuid.New(),
		DashboardID:           dashboardID,
		ReportID:              reportID,
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

	if kpi.Filters == nil {
		kpi.Filters = []Filter{}
	}

	_, err = r.pool.Exec(ctx,
		`INSERT INTO kpis (id, dashboard_id, report_id, data_source_id, label, measurement_name, timeframe, aggregation, filters, comparison_enabled, comparison_display_type, position, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`,
		kpi.ID, kpi.DashboardID, kpi.ReportID, kpi.DataSourceID, kpi.Label, kpi.MeasurementName, kpi.Timeframe, kpi.Aggregation, filtersJSON, kpi.ComparisonEnabled, kpi.ComparisonDisplayType, kpi.Position, kpi.CreatedAt, kpi.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return kpi, nil
}

// GetKPIByID retrieves a KPI by its ID.
func (r *Repository) GetKPIByID(ctx context.Context, id uuid.UUID) (*KPI, error) {
	kpi := &KPI{}
	var filtersJSON []byte
	err := r.pool.QueryRow(ctx,
		`SELECT id, dashboard_id, report_id, data_source_id, label, measurement_name, timeframe, aggregation, filters, comparison_enabled, comparison_display_type, position, created_at, updated_at
		FROM kpis WHERE id = $1`,
		id,
	).Scan(&kpi.ID, &kpi.DashboardID, &kpi.ReportID, &kpi.DataSourceID, &kpi.Label, &kpi.MeasurementName, &kpi.Timeframe, &kpi.Aggregation, &filtersJSON, &kpi.ComparisonEnabled, &kpi.ComparisonDisplayType, &kpi.Position, &kpi.CreatedAt, &kpi.UpdatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(filtersJSON, &kpi.Filters); err != nil {
		kpi.Filters = []Filter{}
	}

	return kpi, nil
}

// GetKPIsByDashboardID retrieves all KPIs for a dashboard.
func (r *Repository) GetKPIsByDashboardID(ctx context.Context, dashboardID uuid.UUID) ([]KPI, error) {
	return r.getKPIs(ctx, "dashboard_id", dashboardID)
}

// GetKPIsByReportID retrieves all KPIs for a report.
func (r *Repository) GetKPIsByReportID(ctx context.Context, reportID uuid.UUID) ([]KPI, error) {
	return r.getKPIs(ctx, "report_id", reportID)
}

func (r *Repository) getKPIs(ctx context.Context, column string, id uuid.UUID) ([]KPI, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, dashboard_id, report_id, data_source_id, label, measurement_name, timeframe, aggregation, filters, comparison_enabled, comparison_display_type, position, created_at, updated_at
		FROM kpis WHERE `+column+` = $1
		ORDER BY position ASC`,
		id,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var kpis []KPI
	for rows.Next() {
		var k KPI
		var filtersJSON []byte
		if err := rows.Scan(&k.ID, &k.DashboardID, &k.ReportID, &k.DataSourceID, &k.Label, &k.MeasurementName, &k.Timeframe, &k.Aggregation, &filtersJSON, &k.ComparisonEnabled, &k.ComparisonDisplayType, &k.Position, &k.CreatedAt, &k.UpdatedAt); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(filtersJSON, &k.Filters); err != nil {
			k.Filters = []Filter{}
		}
		kpis = append(kpis, k)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if kpis == nil {
		kpis = []KPI{}
	}

	return kpis, nil
}

// UpdateKPI updates a KPI's configuration.
func (r *Repository) UpdateKPI(ctx context.Context, id uuid.UUID, req UpdateKPIRequest) error {
	filtersJSON, err := json.Marshal(req.Filters)
	if err != nil {
		return err
	}
	if req.Filters == nil {
		filtersJSON = []byte("[]")
	}

	_, err = r.pool.Exec(ctx,
		`UPDATE kpis SET label = $1, timeframe = $2, aggregation = $3, filters = $4, comparison_enabled = $5, comparison_display_type = $6, updated_at = NOW() WHERE id = $7`,
		req.Label, req.Timeframe, req.Aggregation, filtersJSON, req.ComparisonEnabled, req.ComparisonDisplayType, id,
	)
	return err
}

// DeleteKPI deletes a KPI by its ID.
func (r *Repository) DeleteKPI(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx,
		`DELETE FROM kpis WHERE id = $1`,
		id,
	)
	return err
}

// GetMaxKPIPositionForDashboard gets the maximum position for KPIs in a dashboard.
func (r *Repository) GetMaxKPIPositionForDashboard(ctx context.Context, dashboardID uuid.UUID) (int, error) {
	return r.getMaxKPIPosition(ctx, "dashboard_id", dashboardID)
}

// GetMaxKPIPositionForReport gets the maximum position for KPIs in a report.
func (r *Repository) GetMaxKPIPositionForReport(ctx context.Context, reportID uuid.UUID) (int, error) {
	return r.getMaxKPIPosition(ctx, "report_id", reportID)
}

func (r *Repository) getMaxKPIPosition(ctx context.Context, column string, id uuid.UUID) (int, error) {
	var maxPos *int
	err := r.pool.QueryRow(ctx,
		`SELECT MAX(position) FROM kpis WHERE `+column+` = $1`,
		id,
	).Scan(&maxPos)

	if err != nil {
		return 0, err
	}
	if maxPos == nil {
		return 0, nil
	}

	return *maxPos, nil
}

// UpdateKPIPositionsForDashboard updates the positions of multiple KPIs in a dashboard.
func (r *Repository) UpdateKPIPositionsForDashboard(ctx context.Context, dashboardID uuid.UUID, kpiIDs []uuid.UUID) error {
	return r.updateKPIPositions(ctx, "dashboard_id", dashboardID, kpiIDs)
}

// UpdateKPIPositionsForReport updates the positions of multiple KPIs in a report.
func (r *Repository) UpdateKPIPositionsForReport(ctx context.Context, reportID uuid.UUID, kpiIDs []uuid.UUID) error {
	return r.updateKPIPositions(ctx, "report_id", reportID, kpiIDs)
}

func (r *Repository) updateKPIPositions(ctx context.Context, column string, id uuid.UUID, kpiIDs []uuid.UUID) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	for i, kpiID := range kpiIDs {
		_, err := tx.Exec(ctx,
			`UPDATE kpis SET position = $1, updated_at = NOW() WHERE id = $2 AND `+column+` = $3`,
			i, kpiID, id,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}
