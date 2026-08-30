package artwork

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"time"

	"github.com/maksimovyuriy/artfolio/backend/internal/entity"
	"github.com/maksimovyuriy/artfolio/backend/internal/lib/filestorage"
	"github.com/maksimovyuriy/artfolio/backend/internal/repo"
	"github.com/maksimovyuriy/artfolio/backend/internal/usecase"
)

type UseCase struct {
	repo    repo.ArtworkRepository
	storage filestorage.Artwork
	log     *slog.Logger
}

func NewUseCase(repo repo.ArtworkRepository, storage filestorage.Artwork, log *slog.Logger) *UseCase {
	return &UseCase{repo: repo, storage: storage, log: log}
}

var _ usecase.ArtworkUseCase = (*UseCase)(nil)

func (u *UseCase) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return entity.NewValidationError("id", "must be positive")
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

func (u *UseCase) Reorder(ctx context.Context, artworkIDs []int64) error {
	seen := make(map[int64]struct{}, len(artworkIDs))
	for _, id := range artworkIDs {
		if id <= 0 {
			return entity.NewValidationError("artworkIds", "must contain positive ids")
		}
		if _, exists := seen[id]; exists {
			return entity.NewValidationError("artworkIds", "must not contain duplicates")
		}
		seen[id] = struct{}{}
	}

	if err := u.repo.Reorder(ctx, artworkIDs); errors.Is(err, repo.ErrConflict) {
		return usecase.ErrArtworkOrderConflict
	} else if err != nil {
		return fmt.Errorf("reorder artworks: %w", err)
	}
	return nil
}

func (u *UseCase) Create(ctx context.Context, artwork entity.Artwork, image io.Reader) (entity.Artwork, error) {
	artwork, err := artwork.Validated()
	if err != nil {
		return entity.Artwork{}, err
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
	if artwork.ID <= 0 {
		return entity.NewValidationError("id", "must be positive")
	}
	artwork, err := artwork.Validated()
	if err != nil {
		return err
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
