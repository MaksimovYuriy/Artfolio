package response

import "github.com/maksimovyuriy/artfolio/backend/internal/entity"

type ArtistProfile struct {
	Name            string  `json:"name"`
	Tagline         string  `json:"tagline"`
	Biography       string  `json:"biography"`
	ArtistStatement *string `json:"artistStatement,omitempty"`
	AvatarURL       *string `json:"avatarUrl,omitempty"`
	Email           *string `json:"email,omitempty"`
}

func ArtistProfileFromEntity(profile entity.ArtistProfile) ArtistProfile {
	return ArtistProfile{
		Name:            profile.Name,
		Tagline:         profile.Tagline,
		Biography:       profile.Biography,
		ArtistStatement: profile.ArtistStatement,
		AvatarURL:       profile.AvatarURL,
		Email:           profile.Email,
	}
}
