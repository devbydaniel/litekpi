-- Revert: Rename measurements table back to metrics
ALTER TABLE measurements RENAME TO metrics;

-- Rename indexes back
ALTER INDEX idx_measurements_product_id RENAME TO idx_metrics_product_id;
ALTER INDEX idx_measurements_product_timestamp RENAME TO idx_metrics_product_timestamp;
