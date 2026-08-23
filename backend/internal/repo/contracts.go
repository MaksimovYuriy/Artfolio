package repo

import (
	"context"
	"time"
)

type KeyRepository interface {
	Create(ctx context.Context, keyHash []byte) error
	Find(ctx context.Context, keyHash []byte) (int64, error)
}

type SessionRepository interface {
	Create(
		ctx context.Context,
		adminKeyID int64,
		tokenHash []byte,
		expiresAt time.Time,
	) error

	ExistsActive(
		ctx context.Context,
		tokenHash []byte,
		now time.Time,
	) (bool, error)

	Revoke(
		ctx context.Context,
		tokenHash []byte,
		revokedAt time.Time,
	) error
}
