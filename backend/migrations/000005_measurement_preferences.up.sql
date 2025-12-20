CREATE TABLE measurement_preferences (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    measurement_name VARCHAR(128) NOT NULL,
    preferences JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- One preference per product + measurement combination
    UNIQUE(product_id, measurement_name)
);

-- Index for querying preferences by product
CREATE INDEX idx_measurement_preferences_product_id ON measurement_preferences(product_id);
