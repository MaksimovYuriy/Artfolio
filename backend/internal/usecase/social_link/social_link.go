package sociallink

import (
	"context"
	"fmt"

	"github.com/maksimovyuriy/artfolio/backend/internal/entity"
	"github.com/maksimovyuriy/artfolio/backend/internal/repo"
	"github.com/maksimovyuriy/artfolio/backend/internal/usecase"
)

type UseCase struct {
	profileRepo repo.ArtistProfileRepository
	linkRepo    repo.SocialLinkRepository
}

func NewUseCase(profileRepo repo.ArtistProfileRepository, linkRepo repo.SocialLinkRepository) *UseCase {
	return &UseCase{profileRepo: profileRepo, linkRepo: linkRepo}
}

var _ usecase.SocialLinkUseCase = (*UseCase)(nil)

func (u *UseCase) List(ctx context.Context) ([]entity.SocialLink, error) {
	profile, err := u.profileRepo.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("get artist profile for social links: %w", err)
	}
	links, err := u.linkRepo.List(ctx, profile.ID)
	if err != nil {
		return nil, fmt.Errorf("list social links: %w", err)
	}
	return links, nil
}

func (u *UseCase) Replace(ctx context.Context, links []entity.SocialLink) error {
	normalized := make([]entity.SocialLink, 0, len(links))
	seen := make(map[entity.SocialPlatform]struct{}, len(links))
	for _, link := range links {
		link, err := link.Validated()
		if err != nil {
			return err
		}
		if link.Handle == "" {
			continue
		}
		if _, exists := seen[link.Platform]; exists {
			return entity.NewValidationError("platform", "must not be duplicated")
		}
		seen[link.Platform] = struct{}{}
		normalized = append(normalized, link)
	}

	profile, err := u.profileRepo.Get(ctx)
	if err != nil {
		return fmt.Errorf("get artist profile for social links: %w", err)
	}
	if err := u.linkRepo.Replace(ctx, profile.ID, normalized); err != nil {
		return fmt.Errorf("replace social links: %w", err)
	}
	return nil
}
