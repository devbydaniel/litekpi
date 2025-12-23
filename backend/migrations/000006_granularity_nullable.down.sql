-- Restore granularity as required field
UPDATE metrics SET granularity = 'daily' WHERE granularity IS NULL;
ALTER TABLE metrics ALTER COLUMN granularity SET DEFAULT 'daily';
ALTER TABLE metrics ALTER COLUMN granularity SET NOT NULL;
