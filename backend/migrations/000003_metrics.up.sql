CREATE TABLE metrics (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    name VARCHAR(128) NOT NULL,
    value DOUBLE PRECISION NOT NULL,
    timestamp TIMESTAMPTZ NOT NULL,
    metadata JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    -- Uniqueness constraint for idempotency (product + name + timestamp)
    UNIQUE(product_id, name, timestamp)
);

-- Index for querying metrics by product
CREATE INDEX idx_metrics_product_id ON metrics(product_id);

-- Index for time-range queries
CREATE INDEX idx_metrics_product_timestamp ON metrics(product_id, timestamp);
