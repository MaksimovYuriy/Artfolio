package session

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/maksimovyuriy/artfolio/backend/internal/repo"
)

type Repo struct {
	database *sql.DB
}

func NewRepo(database *sql.DB) *Repo {
	return &Repo{database: database}
}

var _ repo.SessionRepository = (*Repo)(nil)

func (r *Repo) Create(
	ctx context.Context,
	adminKeyID int64,
	tokenHash []byte,
	expiresAt time.Time,
) error {
	const query = `
		INSERT INTO admin_sessions (admin_key_id, token_hash, expires_at)
		VALUES ($1, $2, $3)
	`

	_, err := r.database.ExecContext(ctx, query, adminKeyID, tokenHash, expiresAt)
	if err != nil {
		return fmt.Errorf("create admin session: %w", err)
	}

	return nil
}

func (r *Repo) ExistsActive(
	ctx context.Context,
	tokenHash []byte,
	now time.Time,
) (bool, error) {
	const query = `
		SELECT EXISTS (
			SELECT 1
			FROM admin_sessions
			WHERE token_hash = $1
				AND expires_at > $2
				AND revoked_at IS NULL
		)
	`

	var exists bool
	if err := r.database.QueryRowContext(ctx, query, tokenHash, now).Scan(&exists); err != nil {
		return false, fmt.Errorf("check admin session: %w", err)
	}

	return exists, nil
}

func (r *Repo) Revoke(
	ctx context.Context,
	tokenHash []byte,
	revokedAt time.Time,
) error {
	const query = `
		UPDATE admin_sessions
		SET revoked_at = $2
		WHERE token_hash = $1
			AND revoked_at IS NULL
	`

	_, err := r.database.ExecContext(ctx, query, tokenHash, revokedAt)
	if err != nil {
		return fmt.Errorf("revoke admin session: %w", err)
	}

	return nil
}
