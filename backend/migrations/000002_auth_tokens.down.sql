-- Revert password_hash to NOT NULL (will fail if OAuth-only users exist)
ALTER TABLE users ALTER COLUMN password_hash SET NOT NULL;

-- Drop password reset tokens
DROP INDEX IF EXISTS idx_password_reset_tokens_user_id;
DROP INDEX IF EXISTS idx_password_reset_tokens_token;
DROP TABLE IF EXISTS password_reset_tokens;

-- Drop email verification tokens
DROP INDEX IF EXISTS idx_email_verification_tokens_user_id;
DROP INDEX IF EXISTS idx_email_verification_tokens_token;
DROP TABLE IF EXISTS email_verification_tokens;

-- Drop oauth accounts
DROP TRIGGER IF EXISTS update_oauth_accounts_updated_at ON oauth_accounts;
DROP INDEX IF EXISTS idx_oauth_accounts_provider;
DROP INDEX IF EXISTS idx_oauth_accounts_user_id;
DROP TABLE IF EXISTS oauth_accounts;
