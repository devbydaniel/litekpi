-- Rollback unified metrics migration
-- Note: This will lose any new metrics created after the migration

-- Recreate old tables
CREATE TABLE time_series (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    dashboard_id UUID NOT NULL REFERENCES dashboards(id) ON DELETE CASCADE,
    data_source_id UUID NOT NULL REFERENCES data_sources(id) ON DELETE CASCADE,
    measurement_name VARCHAR(128) NOT NULL,
    title VARCHAR(128),
    date_range VARCHAR(50) NOT NULL DEFAULT 'last_7_days',
    date_from TIMESTAMPTZ,
    date_to TIMESTAMPTZ,
    chart_type VARCHAR(20) NOT NULL DEFAULT 'area',
    split_by VARCHAR(64),
    filters JSONB NOT NULL DEFAULT '[]',
    position INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_time_series_dashboard_id ON time_series(dashboard_id);

CREATE TABLE scalar_metrics (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    dashboard_id UUID NOT NULL REFERENCES dashboards(id) ON DELETE CASCADE,
    label VARCHAR(255) NOT NULL,
    data_source_id UUID NOT NULL REFERENCES data_sources(id) ON DELETE CASCADE,
    measurement_name VARCHAR(128) NOT NULL,
    timeframe VARCHAR(50) NOT NULL DEFAULT 'last_30_days',
    aggregation VARCHAR(20) NOT NULL DEFAULT 'sum',
    filters JSONB NOT NULL DEFAULT '[]',
    comparison_enabled BOOLEAN NOT NULL DEFAULT false,
    comparison_display_type VARCHAR(20),
    position INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_scalar_metrics_dashboard_id ON scalar_metrics(dashboard_id);

-- Migrate data back from metrics to old tables
INSERT INTO scalar_metrics (
    id, dashboard_id, label, position,
    data_source_id, measurement_name, timeframe, filters,
    aggregation, comparison_enabled, comparison_display_type,
    created_at, updated_at
)
SELECT
    id, dashboard_id, label, position,
    data_source_id, measurement_name, timeframe, filters,
    aggregation, comparison_enabled, comparison_display_type,
    created_at, updated_at
FROM metrics
WHERE display_mode = 'scalar';

INSERT INTO time_series (
    id, dashboard_id, title, position,
    data_source_id, measurement_name, date_range, date_from, date_to, filters,
    chart_type, split_by,
    created_at, updated_at
)
SELECT
    id, dashboard_id, label, position,
    data_source_id, measurement_name, timeframe, date_from, date_to, filters,
    chart_type, split_by,
    created_at, updated_at
FROM metrics
WHERE display_mode = 'time_series';

-- Drop unified metrics table
DROP TABLE metrics;
