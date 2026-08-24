package storage

import (
	"context"
	"io"

	"github.com/maksimovyuriy/artfolio/backend/internal/entity"
)

type Artwork interface {
	Save(ctx context.Context, content io.Reader) (entity.StoredArtworkImage, error)
	Delete(ctx context.Context, key string) error
}
