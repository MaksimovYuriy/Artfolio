package artwork

import (
	"context"

	"github.com/maksimovyuriy/artfolio/backend/internal/entity"
	"github.com/maksimovyuriy/artfolio/backend/internal/repo"
	"github.com/maksimovyuriy/artfolio/backend/internal/usecase"
)

type UseCase struct {
	repo repo.ArtworkRepository
}

func NewUseCase(repo repo.ArtworkRepository) *UseCase {
	return &UseCase{repo: repo}
}

var _ usecase.ArtworkUseCase = (*UseCase)(nil)

// Delete implements [usecase.ArtworkUseCase].
func (u *UseCase) Delete(ctx context.Context, id int64) error {
	panic("unimplemented")
}

// ListPublished implements [usecase.ArtworkUseCase].
func (u *UseCase) ListPublished(ctx context.Context) ([]entity.Artwork, error) {
	panic("unimplemented")
}

// ListAll implements [usecase.ArtworkUseCase].
func (u *UseCase) ListAll(ctx context.Context) ([]entity.Artwork, error) {
	panic("unimplemented")
}

// Create implements [usecase.ArtworkUseCase].
func (u *UseCase) Create(ctx context.Context, artwork entity.Artwork) (entity.Artwork, error) {
	panic("unimplemented")
}

// Update implements [usecase.ArtworkUseCase].
func (u *UseCase) Update(ctx context.Context, artwork entity.Artwork) error {
	panic("unimplemented")
}
