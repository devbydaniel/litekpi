-- Domain Model Refactoring
-- Widget → TimeSeries, KPI → ScalarMetric, Delete Reports

-- Drop old tables (fresh start - not live yet)
DROP TABLE IF EXISTS kpis;
DROP TABLE IF EXISTS widgets;
DROP TABLE IF EXISTS reports;

-- Create time_series table (was widgets)
CREATE TABLE time_series (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    dashboard_id UUID NOT NULL REFERENCES dashboards(id) ON DELETE CASCADE,
    title VARCHAR(128),
    position INTEGER NOT NULL DEFAULT 0,
    -- Query fields
    data_source_id UUID NOT NULL REFERENCES data_sources(id) ON DELETE CASCADE,
    measurement_name VARCHAR(128) NOT NULL,
    date_range VARCHAR(50) NOT NULL DEFAULT 'last_30_days',
    date_from TIMESTAMPTZ,
    date_to TIMESTAMPTZ,
    split_by VARCHAR(64),
    filters JSONB NOT NULL DEFAULT '[]',
    -- Display fields
    chart_type VARCHAR(20) NOT NULL DEFAULT 'area' CHECK (chart_type IN ('area', 'bar', 'line')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_time_series_dashboard_id ON time_series(dashboard_id);

-- Create scalar_metrics table (was kpis)
CREATE TABLE scalar_metrics (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    dashboard_id UUID NOT NULL REFERENCES dashboards(id) ON DELETE CASCADE,
    label VARCHAR(255) NOT NULL,
    position INTEGER NOT NULL DEFAULT 0,
    -- Query fields
    data_source_id UUID NOT NULL REFERENCES data_sources(id) ON DELETE CASCADE,
    measurement_name VARCHAR(128) NOT NULL,
    timeframe VARCHAR(50) NOT NULL DEFAULT 'last_30_days' CHECK (timeframe IN ('last_7_days', 'last_30_days', 'this_month', 'last_month')),
    filters JSONB NOT NULL DEFAULT '[]',
    -- Calculation
    aggregation VARCHAR(20) NOT NULL DEFAULT 'sum' CHECK (aggregation IN ('sum', 'average')),
    -- Display fields
    comparison_enabled BOOLEAN NOT NULL DEFAULT false,
    comparison_display_type VARCHAR(20) CHECK (comparison_display_type IN ('percent', 'absolute')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_scalar_metrics_dashboard_id ON scalar_metrics(dashboard_id);
