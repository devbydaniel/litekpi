-- Rollback: Recreate old tables (widgets, kpis, reports)

DROP TABLE IF EXISTS scalar_metrics;
DROP TABLE IF EXISTS time_series;

-- Recreate widgets table
CREATE TABLE widgets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    dashboard_id UUID NOT NULL REFERENCES dashboards(id) ON DELETE CASCADE,
    data_source_id UUID NOT NULL REFERENCES data_sources(id) ON DELETE CASCADE,
    title VARCHAR(128),
    measurement_name VARCHAR(128) NOT NULL,
    chart_type VARCHAR(20) NOT NULL DEFAULT 'area' CHECK (chart_type IN ('area', 'bar', 'line')),
    date_range VARCHAR(50) NOT NULL DEFAULT 'last_30_days',
    date_from TIMESTAMPTZ,
    date_to TIMESTAMPTZ,
    split_by VARCHAR(64),
    filters JSONB NOT NULL DEFAULT '[]',
    position INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_widgets_dashboard_id ON widgets(dashboard_id);

-- Recreate reports table
CREATE TABLE reports (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_reports_organization_id ON reports(organization_id);

-- Recreate kpis table
CREATE TABLE kpis (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    dashboard_id UUID REFERENCES dashboards(id) ON DELETE CASCADE,
    report_id UUID REFERENCES reports(id) ON DELETE CASCADE,
    data_source_id UUID NOT NULL REFERENCES data_sources(id) ON DELETE CASCADE,
    label VARCHAR(255) NOT NULL,
    measurement_name VARCHAR(128) NOT NULL,
    timeframe VARCHAR(50) NOT NULL DEFAULT 'last_30_days',
    aggregation VARCHAR(20) NOT NULL DEFAULT 'sum',
    filters JSONB NOT NULL DEFAULT '[]',
    comparison_enabled BOOLEAN NOT NULL DEFAULT false,
    comparison_display_type VARCHAR(20),
    position INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CHECK (timeframe IN ('last_7_days', 'last_30_days', 'this_month', 'last_month')),
    CHECK (aggregation IN ('sum', 'average')),
    CHECK (comparison_display_type IS NULL OR comparison_display_type IN ('percent', 'absolute')),
    CHECK ((dashboard_id IS NOT NULL AND report_id IS NULL) OR
           (dashboard_id IS NULL AND report_id IS NOT NULL))
);

CREATE INDEX idx_kpis_dashboard_id ON kpis(dashboard_id);
CREATE INDEX idx_kpis_report_id ON kpis(report_id);
