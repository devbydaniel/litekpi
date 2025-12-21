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

// CreateOrganization creates a new organization.
func (r *Repository) CreateOrganization(ctx context.Context, name string) (*Organization, error) {
	org := &Organization{
		ID:        uuid.New(),
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err := r.pool.Exec(ctx,
		`INSERT INTO organizations (id, name, created_at, updated_at)
		VALUES ($1, $2, $3, $4)`,
		org.ID, org.Name, org.CreatedAt, org.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return org, nil
}

// GetOrganizationByID retrieves an organization by ID.
func (r *Repository) GetOrganizationByID(ctx context.Context, id uuid.UUID) (*Organization, error) {
	org := &Organization{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, name, created_at, updated_at FROM organizations WHERE id = $1`,
		id,
	).Scan(&org.ID, &org.Name, &org.CreatedAt, &org.UpdatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return org, nil
}

// CreateUserWithOrg creates a new organization and user in a single transaction.
func (r *Repository) CreateUserWithOrg(ctx context.Context, email, name string, passwordHash *string, orgName string) (*User, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	// Create organization
	org := &Organization{
		ID:        uuid.New(),
		Name:      orgName,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err = tx.Exec(ctx,
		`INSERT INTO organizations (id, name, created_at, updated_at)
		VALUES ($1, $2, $3, $4)`,
		org.ID, org.Name, org.CreatedAt, org.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	// Create user with admin role
	user := &User{
		ID:             uuid.New(),
		Email:          email,
		Name:           name,
		PasswordHash:   passwordHash,
		EmailVerified:  false,
		OrganizationID: org.ID,
		Organization:   org,
		Role:           RoleAdmin,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	_, err = tx.Exec(ctx,
		`INSERT INTO users (id, email, name, password_hash, email_verified, organization_id, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		user.ID, user.Email, user.Name, user.PasswordHash, user.EmailVerified, user.OrganizationID, user.Role, user.CreatedAt, user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return user, nil
}

// CreateUser creates a new user with the specified organization and role.
func (r *Repository) CreateUser(ctx context.Context, email, name string, passwordHash *string, orgID uuid.UUID, role Role) (*User, error) {
	user := &User{
		ID:             uuid.New(),
		Email:          email,
		Name:           name,
		PasswordHash:   passwordHash,
		EmailVerified:  false,
		OrganizationID: orgID,
		Role:           role,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	_, err := r.pool.Exec(ctx,
		`INSERT INTO users (id, email, name, password_hash, email_verified, organization_id, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		user.ID, user.Email, user.Name, user.PasswordHash, user.EmailVerified, user.OrganizationID, user.Role, user.CreatedAt, user.UpdatedAt,
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
		`SELECT id, email, name, password_hash, email_verified, organization_id, role, created_at, updated_at
		FROM users WHERE id = $1`,
		id,
	).Scan(&user.ID, &user.Email, &user.Name, &user.PasswordHash, &user.EmailVerified, &user.OrganizationID, &user.Role, &user.CreatedAt, &user.UpdatedAt)

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
		`SELECT id, email, name, password_hash, email_verified, organization_id, role, created_at, updated_at
		FROM users WHERE email = $1`,
		email,
	).Scan(&user.ID, &user.Email, &user.Name, &user.PasswordHash, &user.EmailVerified, &user.OrganizationID, &user.Role, &user.CreatedAt, &user.UpdatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserWithOrganization retrieves a user with their organization data.
func (r *Repository) GetUserWithOrganization(ctx context.Context, id uuid.UUID) (*User, error) {
	user := &User{Organization: &Organization{}}
	err := r.pool.QueryRow(ctx,
		`SELECT u.id, u.email, u.name, u.password_hash, u.email_verified, u.organization_id, u.role, u.created_at, u.updated_at,
		        o.id, o.name, o.created_at, o.updated_at
		FROM users u
		JOIN organizations o ON u.organization_id = o.id
		WHERE u.id = $1`,
		id,
	).Scan(
		&user.ID, &user.Email, &user.Name, &user.PasswordHash, &user.EmailVerified, &user.OrganizationID, &user.Role, &user.CreatedAt, &user.UpdatedAt,
		&user.Organization.ID, &user.Organization.Name, &user.Organization.CreatedAt, &user.Organization.UpdatedAt,
	)

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

// CreatePendingOAuthAccount creates a new OAuth account without a linked user (pending setup).
func (r *Repository) CreatePendingOAuthAccount(ctx context.Context, provider, providerUserID, providerEmail, providerName string) (*OAuthAccount, error) {
	account := &OAuthAccount{
		ID:             uuid.New(),
		UserID:         nil, // Pending setup
		Provider:       provider,
		ProviderUserID: providerUserID,
		ProviderEmail:  providerEmail,
		ProviderName:   providerName,
		CreatedAt:      time.Now(),
	}

	_, err := r.pool.Exec(ctx,
		`INSERT INTO oauth_accounts (id, user_id, provider, provider_user_id, provider_email, provider_name, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (provider, provider_user_id) DO UPDATE SET
			provider_email = EXCLUDED.provider_email,
			provider_name = EXCLUDED.provider_name`,
		account.ID, account.UserID, account.Provider, account.ProviderUserID, account.ProviderEmail, account.ProviderName, account.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return account, nil
}

// CreateOAuthAccount creates a new OAuth account linked to a user.
func (r *Repository) CreateOAuthAccount(ctx context.Context, account *OAuthAccount) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO oauth_accounts (id, user_id, provider, provider_user_id, provider_email, provider_name, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		account.ID, account.UserID, account.Provider, account.ProviderUserID, account.ProviderEmail, account.ProviderName, account.CreatedAt,
	)
	return err
}

// GetOAuthAccount retrieves an OAuth account by provider and provider user ID.
func (r *Repository) GetOAuthAccount(ctx context.Context, provider, providerUserID string) (*OAuthAccount, error) {
	account := &OAuthAccount{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, user_id, provider, provider_user_id, provider_email, provider_name, created_at
		FROM oauth_accounts WHERE provider = $1 AND provider_user_id = $2`,
		provider, providerUserID,
	).Scan(&account.ID, &account.UserID, &account.Provider, &account.ProviderUserID, &account.ProviderEmail, &account.ProviderName, &account.CreatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return account, nil
}

// GetOAuthAccountByID retrieves an OAuth account by ID.
func (r *Repository) GetOAuthAccountByID(ctx context.Context, id uuid.UUID) (*OAuthAccount, error) {
	account := &OAuthAccount{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, user_id, provider, provider_user_id, provider_email, provider_name, created_at
		FROM oauth_accounts WHERE id = $1`,
		id,
	).Scan(&account.ID, &account.UserID, &account.Provider, &account.ProviderUserID, &account.ProviderEmail, &account.ProviderName, &account.CreatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return account, nil
}

// LinkOAuthAccountToUser links a pending OAuth account to a user.
func (r *Repository) LinkOAuthAccountToUser(ctx context.Context, oauthID, userID uuid.UUID) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE oauth_accounts SET user_id = $1 WHERE id = $2`,
		userID, oauthID,
	)
	return err
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

// CreateInvite creates a new invite.
func (r *Repository) CreateInvite(ctx context.Context, orgID uuid.UUID, email string, role Role, token string, invitedBy uuid.UUID, expiresAt time.Time) (*Invite, error) {
	invite := &Invite{
		ID:             uuid.New(),
		OrganizationID: orgID,
		Email:          email,
		Role:           role,
		Token:          token,
		InvitedBy:      invitedBy,
		ExpiresAt:      expiresAt,
		CreatedAt:      time.Now(),
	}

	_, err := r.pool.Exec(ctx,
		`INSERT INTO invites (id, organization_id, email, role, token, invited_by, expires_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		invite.ID, invite.OrganizationID, invite.Email, invite.Role, invite.Token, invite.InvitedBy, invite.ExpiresAt, invite.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return invite, nil
}

// GetInviteByToken retrieves an invite by its token.
func (r *Repository) GetInviteByToken(ctx context.Context, token string) (*Invite, error) {
	invite := &Invite{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, organization_id, email, role, token, invited_by, expires_at, accepted_at, created_at
		FROM invites WHERE token = $1`,
		token,
	).Scan(&invite.ID, &invite.OrganizationID, &invite.Email, &invite.Role, &invite.Token, &invite.InvitedBy, &invite.ExpiresAt, &invite.AcceptedAt, &invite.CreatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return invite, nil
}

// GetInviteByID retrieves an invite by its ID.
func (r *Repository) GetInviteByID(ctx context.Context, id uuid.UUID) (*Invite, error) {
	invite := &Invite{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, organization_id, email, role, token, invited_by, expires_at, accepted_at, created_at
		FROM invites WHERE id = $1`,
		id,
	).Scan(&invite.ID, &invite.OrganizationID, &invite.Email, &invite.Role, &invite.Token, &invite.InvitedBy, &invite.ExpiresAt, &invite.AcceptedAt, &invite.CreatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return invite, nil
}

// GetPendingInviteByEmail retrieves a pending invite by email within an organization.
func (r *Repository) GetPendingInviteByEmail(ctx context.Context, orgID uuid.UUID, email string) (*Invite, error) {
	invite := &Invite{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, organization_id, email, role, token, invited_by, expires_at, accepted_at, created_at
		FROM invites WHERE organization_id = $1 AND email = $2 AND accepted_at IS NULL`,
		orgID, email,
	).Scan(&invite.ID, &invite.OrganizationID, &invite.Email, &invite.Role, &invite.Token, &invite.InvitedBy, &invite.ExpiresAt, &invite.AcceptedAt, &invite.CreatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return invite, nil
}

// ListPendingInvites retrieves all pending invites for an organization with inviter info.
func (r *Repository) ListPendingInvites(ctx context.Context, orgID uuid.UUID) ([]InviteWithInviter, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT i.id, i.organization_id, i.email, i.role, i.invited_by, i.expires_at, i.accepted_at, i.created_at,
		        u.name, u.email
		FROM invites i
		JOIN users u ON i.invited_by = u.id
		WHERE i.organization_id = $1 AND i.accepted_at IS NULL AND i.expires_at > NOW()
		ORDER BY i.created_at DESC`,
		orgID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var invites []InviteWithInviter
	for rows.Next() {
		var invite InviteWithInviter
		if err := rows.Scan(
			&invite.ID, &invite.OrganizationID, &invite.Email, &invite.Role, &invite.InvitedBy,
			&invite.ExpiresAt, &invite.AcceptedAt, &invite.CreatedAt,
			&invite.InviterName, &invite.InviterEmail,
		); err != nil {
			return nil, err
		}
		invites = append(invites, invite)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return invites, nil
}

// MarkInviteAccepted marks an invite as accepted.
func (r *Repository) MarkInviteAccepted(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE invites SET accepted_at = NOW() WHERE id = $1`,
		id,
	)
	return err
}

// DeleteInvite deletes an invite.
func (r *Repository) DeleteInvite(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx,
		`DELETE FROM invites WHERE id = $1`,
		id,
	)
	return err
}

// ListUsersByOrg retrieves all users for an organization.
func (r *Repository) ListUsersByOrg(ctx context.Context, orgID uuid.UUID) ([]User, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, email, name, password_hash, email_verified, organization_id, role, created_at, updated_at
		FROM users WHERE organization_id = $1
		ORDER BY created_at ASC`,
		orgID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(
			&user.ID, &user.Email, &user.Name, &user.PasswordHash, &user.EmailVerified,
			&user.OrganizationID, &user.Role, &user.CreatedAt, &user.UpdatedAt,
		); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// UpdateUserRole updates a user's role.
func (r *Repository) UpdateUserRole(ctx context.Context, id uuid.UUID, role Role) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE users SET role = $1, updated_at = NOW() WHERE id = $2`,
		role, id,
	)
	return err
}

// DeleteUser deletes a user by ID.
func (r *Repository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx,
		`DELETE FROM users WHERE id = $1`,
		id,
	)
	return err
}

// CountAdmins counts the number of admins in an organization.
func (r *Repository) CountAdmins(ctx context.Context, orgID uuid.UUID) (int, error) {
	var count int
	err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM users WHERE organization_id = $1 AND role = 'admin'`,
		orgID,
	).Scan(&count)
	return count, err
}
