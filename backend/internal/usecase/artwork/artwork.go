package artwork

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/maksimovyuriy/artfolio/backend/internal/entity"
	"github.com/maksimovyuriy/artfolio/backend/internal/lib/storage"
	"github.com/maksimovyuriy/artfolio/backend/internal/repo"
	"github.com/maksimovyuriy/artfolio/backend/internal/usecase"
)

type UseCase struct {
	repo    repo.ArtworkRepository
	storage storage.Artwork
	log     *slog.Logger
}

func NewUseCase(repo repo.ArtworkRepository, storage storage.Artwork, log *slog.Logger) *UseCase {
	return &UseCase{repo: repo, storage: storage, log: log}
}

var _ usecase.ArtworkUseCase = (*UseCase)(nil)

func (u *UseCase) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return usecase.ErrInvalidArtwork
	}

	deleted, err := u.repo.Delete(ctx, id)
	if errors.Is(err, repo.ErrNotFound) {
		return usecase.ErrArtworkNotFound
	}
	if err != nil {
		return fmt.Errorf("delete artwork: %w", err)
	}

	u.deleteImage(ctx, deleted.ImageKey, "delete removed artwork image")
	return nil
}

func (u *UseCase) ListPublished(ctx context.Context) ([]entity.Artwork, error) {
	artworks, err := u.repo.ListPublished(ctx)
	if err != nil {
		return nil, fmt.Errorf("list published artworks: %w", err)
	}
	return artworks, nil
}

func (u *UseCase) ListAll(ctx context.Context) ([]entity.Artwork, error) {
	artworks, err := u.repo.ListAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("list all artworks: %w", err)
	}
	return artworks, nil
}

func (u *UseCase) Create(ctx context.Context, artwork entity.Artwork, image io.Reader) (entity.Artwork, error) {
	artwork = normalize(artwork)
	if !valid(artwork, false) {
		return entity.Artwork{}, usecase.ErrInvalidArtwork
	}
	if image == nil {
		return entity.Artwork{}, usecase.ErrArtworkImageRequired
	}

	stored, err := u.storage.Save(ctx, image)
	if err != nil {
		return entity.Artwork{}, fmt.Errorf("save artwork image: %w", err)
	}
	artwork.ID = 0
	artwork.ImageKey = stored.Key
	artwork.CreatedAt = time.Time{}
	artwork.UpdatedAt = time.Time{}

	created, err := u.repo.Create(ctx, artwork)
	if err != nil {
		u.deleteImage(ctx, stored.Key, "clean up image after artwork creation failure")
		return entity.Artwork{}, fmt.Errorf("create artwork: %w", err)
	}
	return created, nil
}

func (u *UseCase) Update(ctx context.Context, artwork entity.Artwork, image io.Reader) error {
	artwork = normalize(artwork)
	if !valid(artwork, true) {
		return usecase.ErrInvalidArtwork
	}

	if image == nil {
		_, err := u.repo.Update(ctx, artwork)
		return updateError(err)
	}

	stored, err := u.storage.Save(ctx, image)
	if err != nil {
		return fmt.Errorf("save updated artwork image: %w", err)
	}
	artwork.ImageKey = stored.Key

	previous, err := u.repo.UpdateWithImage(ctx, artwork)
	if err != nil {
		u.deleteImage(ctx, stored.Key, "clean up image after artwork update failure")
		return updateError(err)
	}

	u.deleteImage(ctx, previous.ImageKey, "delete replaced artwork image")
	return nil
}

func updateError(err error) error {
	if errors.Is(err, repo.ErrNotFound) {
		return usecase.ErrArtworkNotFound
	}
	if err != nil {
		return fmt.Errorf("update artwork: %w", err)
	}
	return nil
}

func normalize(artwork entity.Artwork) entity.Artwork {
	artwork.Title = strings.TrimSpace(artwork.Title)
	artwork.Description = strings.TrimSpace(artwork.Description)
	artwork.Technique = strings.TrimSpace(artwork.Technique)
	artwork.ImageAlt = strings.TrimSpace(artwork.ImageAlt)
	return artwork
}

func valid(artwork entity.Artwork, requireID bool) bool {
	if requireID && artwork.ID <= 0 {
		return false
	}
	if artwork.Title == "" || utf8.RuneCountInString(artwork.Title) > 256 {
		return false
	}
	if utf8.RuneCountInString(artwork.Technique) > 256 || utf8.RuneCountInString(artwork.ImageAlt) > 256 {
		return false
	}
	if artwork.Year != nil && (*artwork.Year < 0 || *artwork.Year > 9999) {
		return false
	}
	return artwork.Position >= 0
}

func (u *UseCase) deleteImage(ctx context.Context, key, message string) {
	if key == "" {
		return
	}
	cleanupCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 5*time.Second)
	defer cancel()
	if err := u.storage.Delete(cleanupCtx, key); err != nil {
		u.log.Error(message, slog.String("image_key", key), slog.Any("error", err))
	}
}
