package sociallink

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"unicode/utf8"

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
		if !validPlatform(link.Platform) {
			return fmt.Errorf("%w: unsupported platform %s", usecase.ErrInvalidSocialLinks, link.Platform)
		}
		link.Handle = normalizeHandle(link.Platform, link.Handle)
		if link.Handle == "" {
			continue
		}
		if !validHandle(link.Handle) {
			return fmt.Errorf("%w: invalid %s handle", usecase.ErrInvalidSocialLinks, link.Platform)
		}
		if _, exists := seen[link.Platform]; exists {
			return fmt.Errorf("%w: duplicate platform %s", usecase.ErrInvalidSocialLinks, link.Platform)
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

func validPlatform(platform entity.SocialPlatform) bool {
	switch platform {
	case entity.SocialPlatformTelegram,
		entity.SocialPlatformInstagram,
		entity.SocialPlatformVK,
		entity.SocialPlatformBehance:
		return true
	default:
		return false
	}
}

func validHandle(handle string) bool {
	if utf8.RuneCountInString(handle) > 256 {
		return false
	}
	for _, character := range handle {
		if character >= 'a' && character <= 'z' ||
			character >= 'A' && character <= 'Z' ||
			character >= '0' && character <= '9' ||
			character == '_' || character == '-' || character == '.' {
			continue
		}
		return false
	}
	return handle != ""
}

func normalizeHandle(platform entity.SocialPlatform, handle string) string {
	handle = strings.TrimSpace(handle)
	if handle == "" {
		return ""
	}

	if parsed, err := url.Parse(handle); err == nil && parsed.Host != "" {
		host := strings.TrimPrefix(strings.ToLower(parsed.Hostname()), "www.")
		if host == platformHost(platform) && parsed.RawQuery == "" && parsed.Fragment == "" {
			handle = strings.Trim(parsed.Path, "/")
		}
	}

	return strings.TrimPrefix(strings.TrimSpace(handle), "@")
}

func platformHost(platform entity.SocialPlatform) string {
	switch platform {
	case entity.SocialPlatformTelegram:
		return "t.me"
	case entity.SocialPlatformInstagram:
		return "instagram.com"
	case entity.SocialPlatformVK:
		return "vk.com"
	case entity.SocialPlatformBehance:
		return "behance.net"
	default:
		return ""
	}
}
