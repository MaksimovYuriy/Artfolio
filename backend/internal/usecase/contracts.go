package usecase

import (
	"context"
	"errors"

	"github.com/maksimovyuriy/artfolio/backend/internal/entity"
)

var ErrInvalidArtistProfile = errors.New("invalid artist profile")

type KeyCreator interface {
	Create(ctx context.Context) (string, error)
}

type SessionUseCase interface {
	Create(ctx context.Context, accessKey string) (entity.Session, error)
	Verify(ctx context.Context, token string) (bool, error)
	Revoke(ctx context.Context, token string) error
}

type ArtistProfileUseCase interface {
	Get(ctx context.Context) (entity.ArtistProfile, error)
	Update(ctx context.Context, profile entity.ArtistProfile) error
}

type ArtworkUseCase interface {
	ListPublished(ctx context.Context) ([]entity.Artwork, error)
	ListAll(ctx context.Context) ([]entity.Artwork, error)
	Create(ctx context.Context, artwork entity.Artwork) (entity.Artwork, error)
	Update(ctx context.Context, artwork entity.Artwork) error
	Delete(ctx context.Context, id int64) error
}
