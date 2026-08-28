package session

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/maksimovyuriy/artfolio/backend/internal/entity"
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
) (int64, error) {
	const query = `
		INSERT INTO admin_sessions (admin_key_id, token_hash, expires_at)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	var sessionID int64
	if err := r.database.QueryRowContext(ctx, query, adminKeyID, tokenHash, expiresAt).Scan(&sessionID); err != nil {
		return 0, fmt.Errorf("create admin session: %w", err)
	}

	return sessionID, nil
}

func (r *Repo) FindActive(
	ctx context.Context,
	tokenHash []byte,
	now time.Time,
) (entity.AuthenticatedSession, error) {
	const query = `
		SELECT id, admin_key_id
		FROM admin_sessions
		WHERE token_hash = $1
			AND expires_at > $2
			AND revoked_at IS NULL
	`

	var session entity.AuthenticatedSession
	if err := r.database.QueryRowContext(ctx, query, tokenHash, now).Scan(&session.ID, &session.ActorID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.AuthenticatedSession{}, repo.ErrNotFound
		}
		return entity.AuthenticatedSession{}, fmt.Errorf("find active admin session: %w", err)
	}

	return session, nil
}

func (r *Repo) Revoke(
	ctx context.Context,
	tokenHash []byte,
	revokedAt time.Time,
) (entity.AuthenticatedSession, error) {
	const query = `
		UPDATE admin_sessions
		SET revoked_at = $2
		WHERE token_hash = $1
			AND revoked_at IS NULL
		RETURNING id, admin_key_id
	`

	var session entity.AuthenticatedSession
	if err := r.database.QueryRowContext(ctx, query, tokenHash, revokedAt).Scan(&session.ID, &session.ActorID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.AuthenticatedSession{}, repo.ErrNotFound
		}
		return entity.AuthenticatedSession{}, fmt.Errorf("revoke admin session: %w", err)
	}

	return session, nil
}
