-- Make granularity nullable (only required for time_series display mode)
ALTER TABLE metrics ALTER COLUMN granularity DROP NOT NULL;
ALTER TABLE metrics ALTER COLUMN granularity DROP DEFAULT;

-- Set existing scalar metrics to NULL
UPDATE metrics SET granularity = NULL WHERE display_mode = 'scalar';
