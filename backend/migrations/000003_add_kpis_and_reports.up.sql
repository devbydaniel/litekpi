-- Reports table
CREATE TABLE reports (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_reports_organization_id ON reports(organization_id);

CREATE TRIGGER update_reports_updated_at
    BEFORE UPDATE ON reports
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- KPIs table (can belong to dashboard OR report)
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
CREATE INDEX idx_kpis_data_source_id ON kpis(data_source_id);

CREATE TRIGGER update_kpis_updated_at
    BEFORE UPDATE ON kpis
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
