package storage

import (
	"context"
	"errors"
	"io"

	"github.com/maksimovyuriy/artfolio/backend/internal/entity"
)

var (
	ErrFileTooLarge       = errors.New("artwork image is too large")
	ErrInvalidImage       = errors.New("invalid artwork image")
	ErrImageTooManyPixels = errors.New("artwork image has too many pixels")
	ErrInvalidKey         = errors.New("invalid artwork image key")
)

type Artwork interface {
	Save(ctx context.Context, content io.Reader) (entity.StoredArtworkImage, error)
	Delete(ctx context.Context, key string) error
}
