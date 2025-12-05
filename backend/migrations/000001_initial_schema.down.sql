-- Rollback initial schema

-- Drop triggers first
DROP TRIGGER IF EXISTS update_products_updated_at ON products;
DROP TRIGGER IF EXISTS update_users_updated_at ON users;

-- Drop trigger function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop tables in reverse order (respecting foreign key constraints)
DROP TABLE IF EXISTS data_points;
DROP TABLE IF EXISTS products;
DROP TABLE IF EXISTS users;

-- Drop extensions
DROP EXTENSION IF EXISTS "uuid-ossp";
