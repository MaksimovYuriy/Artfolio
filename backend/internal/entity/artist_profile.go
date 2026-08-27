package entity

import (
	"net/mail"
	"strings"
)

type ArtistProfile struct {
	ID              int64
	Name            string
	Tagline         string
	Biography       string
	ArtistStatement *string
	Email           *string
	SocialLinks     []SocialLink
}

func (profile ArtistProfile) Validated() (ArtistProfile, error) {
	profile.normalize()

	if profile.Name == "" {
		return ArtistProfile{}, NewValidationError("name", "is required")
	}
	if profile.Email != nil {
		address, err := mail.ParseAddress(*profile.Email)
		if err != nil || address.Address != *profile.Email {
			return ArtistProfile{}, NewValidationError("email", "is invalid")
		}
	}
	return profile, nil
}

func (profile *ArtistProfile) normalize() {
	profile.Name = strings.TrimSpace(profile.Name)
	profile.Tagline = strings.TrimSpace(profile.Tagline)
	profile.Biography = strings.TrimSpace(profile.Biography)
	profile.ArtistStatement = normalizeOptionalString(profile.ArtistStatement)
	profile.Email = normalizeOptionalString(profile.Email)
}

func normalizeOptionalString(value *string) *string {
	if value == nil {
		return nil
	}

	normalized := strings.TrimSpace(*value)
	if normalized == "" {
		return nil
	}

	return &normalized
}
