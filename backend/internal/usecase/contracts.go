package usecase

import (
	"context"
	"errors"
	"io"

	"github.com/maksimovyuriy/artfolio/backend/internal/entity"
)

var (
	ErrArtworkNotFound      = errors.New("artwork not found")
	ErrArtworkImageRequired = errors.New("artwork image is required")
	ErrArtworkOrderConflict = errors.New("artwork order conflict")
	ErrInvalidSession       = errors.New("invalid session")
)

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
	Create(ctx context.Context, artwork entity.Artwork, image io.Reader) (entity.Artwork, error)
	Update(ctx context.Context, artwork entity.Artwork, image io.Reader) error
	Reorder(ctx context.Context, artworkIDs []int64) error
	Delete(ctx context.Context, id int64) error
}

type SocialLinkUseCase interface {
	List(ctx context.Context) ([]entity.SocialLink, error)
	Replace(ctx context.Context, links []entity.SocialLink) error
}
