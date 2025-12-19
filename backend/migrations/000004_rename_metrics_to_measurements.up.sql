-- Rename metrics table to measurements
ALTER TABLE metrics RENAME TO measurements;

-- Rename indexes
ALTER INDEX idx_metrics_product_id RENAME TO idx_measurements_product_id;
ALTER INDEX idx_metrics_product_timestamp RENAME TO idx_measurements_product_timestamp;
