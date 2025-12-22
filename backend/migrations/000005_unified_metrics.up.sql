-- Unified Metrics
-- Combines scalar_metrics and time_series into a single metrics table

-- Create unified metrics table
CREATE TABLE metrics (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    dashboard_id UUID NOT NULL REFERENCES dashboards(id) ON DELETE CASCADE,
    label VARCHAR(255) NOT NULL,
    position INTEGER NOT NULL DEFAULT 0,

    -- Query fields
    data_source_id UUID NOT NULL REFERENCES data_sources(id) ON DELETE CASCADE,
    measurement_name VARCHAR(128) NOT NULL,
    timeframe VARCHAR(50) NOT NULL DEFAULT 'last_30_days',
    date_from TIMESTAMPTZ,
    date_to TIMESTAMPTZ,
    filters JSONB NOT NULL DEFAULT '[]',

    -- Aggregation
    aggregation VARCHAR(20) NOT NULL DEFAULT 'sum'
        CHECK (aggregation IN ('sum', 'average', 'count', 'count_unique')),
    aggregation_key VARCHAR(64), -- For count_unique: which metadata key to count distinct values of
    granularity VARCHAR(20) NOT NULL DEFAULT 'daily'
        CHECK (granularity IN ('daily', 'weekly', 'monthly')),

    -- Display mode
    display_mode VARCHAR(20) NOT NULL DEFAULT 'scalar'
        CHECK (display_mode IN ('scalar', 'time_series')),

    -- Scalar display options
    comparison_enabled BOOLEAN NOT NULL DEFAULT false,
    comparison_display_type VARCHAR(20)
        CHECK (comparison_display_type IN ('percent', 'absolute')),

    -- Time series display options
    chart_type VARCHAR(20) CHECK (chart_type IN ('area', 'bar', 'line')),
    split_by VARCHAR(64),

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Constraint: aggregation_key required for count_unique
    CONSTRAINT check_aggregation_key_required
        CHECK (aggregation != 'count_unique' OR aggregation_key IS NOT NULL),
    -- Constraint: chart_type required for time_series
    CONSTRAINT check_chart_type_for_time_series
        CHECK (display_mode != 'time_series' OR chart_type IS NOT NULL)
);

CREATE INDEX idx_metrics_dashboard_id ON metrics(dashboard_id);
CREATE INDEX idx_metrics_data_source_id ON metrics(data_source_id);

-- Migrate scalar_metrics to metrics
INSERT INTO metrics (
    id, dashboard_id, label, position,
    data_source_id, measurement_name, timeframe, filters,
    aggregation, granularity, display_mode,
    comparison_enabled, comparison_display_type,
    created_at, updated_at
)
SELECT
    id, dashboard_id, label, position,
    data_source_id, measurement_name, timeframe, filters,
    aggregation, 'daily', 'scalar',
    comparison_enabled, comparison_display_type,
    created_at, updated_at
FROM scalar_metrics;

-- Migrate time_series to metrics
-- Note: time_series.title becomes metrics.label
INSERT INTO metrics (
    id, dashboard_id, label, position,
    data_source_id, measurement_name, timeframe, date_from, date_to, filters,
    aggregation, granularity, display_mode,
    chart_type, split_by,
    created_at, updated_at
)
SELECT
    id, dashboard_id, COALESCE(title, measurement_name), position,
    data_source_id, measurement_name, date_range, date_from, date_to, filters,
    'sum', 'daily', 'time_series',
    chart_type, split_by,
    created_at, updated_at
FROM time_series;

-- Drop old tables
DROP TABLE scalar_metrics;
DROP TABLE time_series;
