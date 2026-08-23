package artistprofile

import (
	"context"
	"fmt"
	"net/mail"
	"strings"
	"unicode/utf8"

	"github.com/maksimovyuriy/artfolio/backend/internal/entity"
	"github.com/maksimovyuriy/artfolio/backend/internal/repo"
	"github.com/maksimovyuriy/artfolio/backend/internal/usecase"
)

type UseCase struct {
	repo repo.ArtistProfileRepository
}

func NewUseCase(repo repo.ArtistProfileRepository) *UseCase {
	return &UseCase{repo: repo}
}

var _ usecase.ArtistProfileUseCase = (*UseCase)(nil)

func (u *UseCase) Get(ctx context.Context) (entity.ArtistProfile, error) {
	return u.repo.Get(ctx)
}

func (u *UseCase) Update(ctx context.Context, profile entity.ArtistProfile) error {
	profile = normalize(profile)
	if err := validate(profile); err != nil {
		return err
	}

	return u.repo.Update(ctx, profile)
}

func normalize(profile entity.ArtistProfile) entity.ArtistProfile {
	profile.Name = strings.TrimSpace(profile.Name)
	profile.Tagline = strings.TrimSpace(profile.Tagline)
	profile.Biography = strings.TrimSpace(profile.Biography)
	profile.ArtistStatement = normalizeOptional(profile.ArtistStatement)
	profile.AvatarURL = normalizeOptional(profile.AvatarURL)
	profile.Email = normalizeOptional(profile.Email)
	return profile
}

func normalizeOptional(value *string) *string {
	if value == nil {
		return nil
	}

	normalized := strings.TrimSpace(*value)
	if normalized == "" {
		return nil
	}

	return &normalized
}

func validate(profile entity.ArtistProfile) error {
	if profile.Name == "" {
		return fmt.Errorf("%w: name is required", usecase.ErrInvalidArtistProfile)
	}
	if utf8.RuneCountInString(profile.Name) > 64 {
		return fmt.Errorf("%w: name must not exceed 64 characters", usecase.ErrInvalidArtistProfile)
	}
	if utf8.RuneCountInString(profile.Tagline) > 256 {
		return fmt.Errorf("%w: tagline must not exceed 256 characters", usecase.ErrInvalidArtistProfile)
	}
	if profile.Email != nil {
		address, err := mail.ParseAddress(*profile.Email)
		if err != nil || address.Address != *profile.Email {
			return fmt.Errorf("%w: email is invalid", usecase.ErrInvalidArtistProfile)
		}
	}

	return nil
}
