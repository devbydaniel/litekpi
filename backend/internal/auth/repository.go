package auth

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository handles database operations for authentication.
type Repository struct {
	pool *pgxpool.Pool
}

// NewRepository creates a new auth repository.
func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

// CreateUser creates a new user with email and password hash.
func (r *Repository) CreateUser(ctx context.Context, email string, passwordHash *string) (*User, error) {
	user := &User{
		ID:            uuid.New(),
		Email:         email,
		PasswordHash:  passwordHash,
		EmailVerified: false,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	_, err := r.pool.Exec(ctx,
		`INSERT INTO users (id, email, password_hash, email_verified, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		user.ID, user.Email, user.PasswordHash, user.EmailVerified, user.CreatedAt, user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserByID retrieves a user by their ID.
func (r *Repository) GetUserByID(ctx context.Context, id uuid.UUID) (*User, error) {
	user := &User{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, email, password_hash, email_verified, created_at, updated_at
		FROM users WHERE id = $1`,
		id,
	).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.EmailVerified, &user.CreatedAt, &user.UpdatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserByEmail retrieves a user by their email.
func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	user := &User{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, email, password_hash, email_verified, created_at, updated_at
		FROM users WHERE email = $1`,
		email,
	).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.EmailVerified, &user.CreatedAt, &user.UpdatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

// UpdateUserEmailVerified updates the email_verified status for a user.
func (r *Repository) UpdateUserEmailVerified(ctx context.Context, id uuid.UUID, verified bool) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE users SET email_verified = $1 WHERE id = $2`,
		verified, id,
	)
	return err
}

// UpdateUserPassword updates the password hash for a user.
func (r *Repository) UpdateUserPassword(ctx context.Context, id uuid.UUID, passwordHash string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE users SET password_hash = $1 WHERE id = $2`,
		passwordHash, id,
	)
	return err
}

// CreateOAuthAccount creates a new OAuth account linked to a user.
func (r *Repository) CreateOAuthAccount(ctx context.Context, account *OAuthAccount) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO oauth_accounts (id, user_id, provider, provider_user_id, provider_email, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		account.ID, account.UserID, account.Provider, account.ProviderUserID, account.ProviderEmail, account.CreatedAt, account.UpdatedAt,
	)
	return err
}

// GetOAuthAccount retrieves an OAuth account by provider and provider user ID.
func (r *Repository) GetOAuthAccount(ctx context.Context, provider, providerUserID string) (*OAuthAccount, error) {
	account := &OAuthAccount{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, user_id, provider, provider_user_id, provider_email, created_at, updated_at
		FROM oauth_accounts WHERE provider = $1 AND provider_user_id = $2`,
		provider, providerUserID,
	).Scan(&account.ID, &account.UserID, &account.Provider, &account.ProviderUserID, &account.ProviderEmail, &account.CreatedAt, &account.UpdatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return account, nil
}

// CreateEmailVerificationToken creates a new email verification token.
func (r *Repository) CreateEmailVerificationToken(ctx context.Context, userID uuid.UUID, token string, expiresAt time.Time) error {
	id := uuid.New()
	_, err := r.pool.Exec(ctx,
		`INSERT INTO email_verification_tokens (id, user_id, token, expires_at, created_at)
		VALUES ($1, $2, $3, $4, $5)`,
		id, userID, token, expiresAt, time.Now(),
	)
	return err
}

// GetEmailVerificationToken retrieves an email verification token by token string.
func (r *Repository) GetEmailVerificationToken(ctx context.Context, token string) (*EmailVerificationToken, error) {
	evt := &EmailVerificationToken{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, user_id, token, expires_at, created_at
		FROM email_verification_tokens WHERE token = $1`,
		token,
	).Scan(&evt.ID, &evt.UserID, &evt.Token, &evt.ExpiresAt, &evt.CreatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return evt, nil
}

// DeleteEmailVerificationTokensByUserID deletes all email verification tokens for a user.
func (r *Repository) DeleteEmailVerificationTokensByUserID(ctx context.Context, userID uuid.UUID) error {
	_, err := r.pool.Exec(ctx,
		`DELETE FROM email_verification_tokens WHERE user_id = $1`,
		userID,
	)
	return err
}

// CreatePasswordResetToken creates a new password reset token.
func (r *Repository) CreatePasswordResetToken(ctx context.Context, userID uuid.UUID, token string, expiresAt time.Time) error {
	id := uuid.New()
	_, err := r.pool.Exec(ctx,
		`INSERT INTO password_reset_tokens (id, user_id, token, expires_at, created_at)
		VALUES ($1, $2, $3, $4, $5)`,
		id, userID, token, expiresAt, time.Now(),
	)
	return err
}

// GetPasswordResetToken retrieves a password reset token by token string.
func (r *Repository) GetPasswordResetToken(ctx context.Context, token string) (*PasswordResetToken, error) {
	prt := &PasswordResetToken{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, user_id, token, expires_at, used, created_at
		FROM password_reset_tokens WHERE token = $1`,
		token,
	).Scan(&prt.ID, &prt.UserID, &prt.Token, &prt.ExpiresAt, &prt.Used, &prt.CreatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return prt, nil
}

// MarkPasswordResetTokenUsed marks a password reset token as used.
func (r *Repository) MarkPasswordResetTokenUsed(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE password_reset_tokens SET used = true WHERE id = $1`,
		id,
	)
	return err
}

// DeletePasswordResetTokensByUserID deletes all password reset tokens for a user.
func (r *Repository) DeletePasswordResetTokensByUserID(ctx context.Context, userID uuid.UUID) error {
	_, err := r.pool.Exec(ctx,
		`DELETE FROM password_reset_tokens WHERE user_id = $1`,
		userID,
	)
	return err
}
