package repo

import (
	"context"
	"time"

	"github.com/maksimovyuriy/artfolio/backend/internal/entity"
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

type ArtistProfileRepository interface {
	Get(ctx context.Context) (entity.ArtistProfile, error)
	Update(ctx context.Context, profile entity.ArtistProfile) error
}

type ArtworkRepository interface {
	ListPublished(ctx context.Context) ([]entity.Artwork, error)
	ListAll(ctx context.Context) ([]entity.Artwork, error)
	Create(ctx context.Context, artwork entity.Artwork) (entity.Artwork, error)
	Update(ctx context.Context, artwork entity.Artwork) error
	Delete(ctx context.Context, id int64) error
}
