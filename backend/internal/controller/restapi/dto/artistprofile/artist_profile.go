package artistprofile

import "github.com/maksimovyuriy/artfolio/backend/internal/entity"

type UpdateArtistProfileRequest struct {
	Name            string  `json:"name"`
	Tagline         string  `json:"tagline"`
	Biography       string  `json:"biography"`
	ArtistStatement *string `json:"artistStatement"`
	AvatarURL       *string `json:"avatarUrl"`
	Email           *string `json:"email"`
}

func (r UpdateArtistProfileRequest) ArtistProfile() entity.ArtistProfile {
	return entity.ArtistProfile{
		Name:            r.Name,
		Tagline:         r.Tagline,
		Biography:       r.Biography,
		ArtistStatement: r.ArtistStatement,
		AvatarURL:       r.AvatarURL,
		Email:           r.Email,
	}
}

type ArtistProfileResponse struct {
	Name            string  `json:"name"`
	Tagline         string  `json:"tagline"`
	Biography       string  `json:"biography"`
	ArtistStatement *string `json:"artistStatement,omitempty"`
	AvatarURL       *string `json:"avatarUrl,omitempty"`
	Email           *string `json:"email,omitempty"`
}

func ArtistProfileResponseFromEntity(profile entity.ArtistProfile) ArtistProfileResponse {
	return ArtistProfileResponse{
		Name:            profile.Name,
		Tagline:         profile.Tagline,
		Biography:       profile.Biography,
		ArtistStatement: profile.ArtistStatement,
		AvatarURL:       profile.AvatarURL,
		Email:           profile.Email,
	}
}
