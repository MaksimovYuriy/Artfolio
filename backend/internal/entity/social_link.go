package entity

import (
	"net/url"
	"strings"
)

type SocialPlatform string

const (
	SocialPlatformTelegram  SocialPlatform = "telegram"
	SocialPlatformInstagram SocialPlatform = "instagram"
	SocialPlatformVK        SocialPlatform = "vk"
	SocialPlatformBehance   SocialPlatform = "behance"
)

type SocialLink struct {
	ArtistProfileID int64
	Platform        SocialPlatform
	Handle          string
}

func (link SocialLink) Validated() (SocialLink, error) {
	if !link.Platform.Valid() {
		return SocialLink{}, NewValidationError("platform", "is unsupported")
	}
	link.normalize()
	if link.Handle == "" {
		return link, nil
	}
	if !validSocialHandle(link.Handle) {
		return SocialLink{}, NewValidationError("handle", "is invalid")
	}
	return link, nil
}

func (link *SocialLink) normalize() {
	handle := strings.TrimSpace(link.Handle)
	if handle == "" {
		link.Handle = ""
		return
	}

	if parsed, err := url.Parse(handle); err == nil && parsed.Host != "" {
		host := strings.TrimPrefix(strings.ToLower(parsed.Hostname()), "www.")
		if host == link.Platform.host() && parsed.RawQuery == "" && parsed.Fragment == "" {
			handle = strings.Trim(parsed.Path, "/")
		}
	}

	link.Handle = strings.TrimPrefix(strings.TrimSpace(handle), "@")
}

func (platform SocialPlatform) Valid() bool {
	switch platform {
	case SocialPlatformTelegram, SocialPlatformInstagram, SocialPlatformVK, SocialPlatformBehance:
		return true
	default:
		return false
	}
}

func (platform SocialPlatform) host() string {
	switch platform {
	case SocialPlatformTelegram:
		return "t.me"
	case SocialPlatformInstagram:
		return "instagram.com"
	case SocialPlatformVK:
		return "vk.com"
	case SocialPlatformBehance:
		return "behance.net"
	default:
		return ""
	}
}

func validSocialHandle(handle string) bool {
	if handle == "" {
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
	return true
}
