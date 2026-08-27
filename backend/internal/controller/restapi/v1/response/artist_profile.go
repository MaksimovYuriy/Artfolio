package response

import "github.com/maksimovyuriy/artfolio/backend/internal/entity"

type ArtistProfile struct {
	Name            string       `json:"name"`
	Tagline         string       `json:"tagline"`
	Biography       string       `json:"biography"`
	ArtistStatement *string      `json:"artistStatement,omitempty"`
	Email           *string      `json:"email,omitempty"`
	SocialLinks     []SocialLink `json:"socialLinks"`
}

func ArtistProfileFromEntity(profile entity.ArtistProfile) ArtistProfile {
	return ArtistProfile{
		Name:            profile.Name,
		Tagline:         profile.Tagline,
		Biography:       profile.Biography,
		ArtistStatement: profile.ArtistStatement,
		Email:           profile.Email,
		SocialLinks:     SocialLinksFromEntities(profile.SocialLinks),
	}
}
