package artistprofile

import (
	"context"
	"fmt"

	"github.com/maksimovyuriy/artfolio/backend/internal/entity"
	"github.com/maksimovyuriy/artfolio/backend/internal/repo"
	"github.com/maksimovyuriy/artfolio/backend/internal/usecase"
)

type UseCase struct {
	profileRepo    repo.ArtistProfileRepository
	socialLinkRepo repo.SocialLinkRepository
}

func NewUseCase(profileRepo repo.ArtistProfileRepository, socialLinkRepo repo.SocialLinkRepository) *UseCase {
	return &UseCase{profileRepo: profileRepo, socialLinkRepo: socialLinkRepo}
}

var _ usecase.ArtistProfileUseCase = (*UseCase)(nil)

func (u *UseCase) Get(ctx context.Context) (entity.ArtistProfile, error) {
	profile, err := u.profileRepo.Get(ctx)
	if err != nil {
		return entity.ArtistProfile{}, err
	}
	links, err := u.socialLinkRepo.List(ctx, profile.ID)
	if err != nil {
		return entity.ArtistProfile{}, fmt.Errorf("list artist social links: %w", err)
	}
	profile.SocialLinks = links
	return profile, nil
}

func (u *UseCase) Update(ctx context.Context, profile entity.ArtistProfile) error {
	profile, err := profile.Validated()
	if err != nil {
		return err
	}

	return u.profileRepo.Update(ctx, profile)
}
