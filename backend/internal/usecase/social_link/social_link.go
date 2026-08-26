package sociallink

import (
	"context"

	"github.com/maksimovyuriy/artfolio/backend/internal/entity"
	"github.com/maksimovyuriy/artfolio/backend/internal/repo"
	"github.com/maksimovyuriy/artfolio/backend/internal/usecase"
)

type UseCase struct {
	repo repo.SocialLinkRepository
}

func NewUseCase(repo repo.SocialLinkRepository) *UseCase {
	return &UseCase{repo: repo}
}

var _ usecase.SocialLinkUseCase = (*UseCase)(nil)

// Create implements [usecase.SocialLinkUseCase].
func (u *UseCase) Create(ctx context.Context, link entity.SocialLink) (entity.SocialLink, error) {
	panic("unimplemented")
}

// Delete implements [usecase.SocialLinkUseCase].
func (u *UseCase) Delete(ctx context.Context, id int64) (entity.SocialLink, error) {
	panic("unimplemented")
}

// List implements [usecase.SocialLinkUseCase].
func (u *UseCase) List(ctx context.Context, artistId int64) ([]entity.SocialLink, error) {
	panic("unimplemented")
}

// Update implements [usecase.SocialLinkUseCase].
func (u *UseCase) Update(ctx context.Context, link entity.SocialLink) (entity.SocialLink, error) {
	panic("unimplemented")
}

func validateHandle(handle string) bool {
	return false
}

func normalize(handle string) string {
	return handle
}
