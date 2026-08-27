package request

import "github.com/maksimovyuriy/artfolio/backend/internal/entity"

type UpdateArtistProfile struct {
	Name            string  `json:"name"`
	Tagline         string  `json:"tagline"`
	Biography       string  `json:"biography"`
	ArtistStatement *string `json:"artistStatement"`
	Email           *string `json:"email"`
}

func (r UpdateArtistProfile) ArtistProfile() entity.ArtistProfile {
	return entity.ArtistProfile{
		Name:            r.Name,
		Tagline:         r.Tagline,
		Biography:       r.Biography,
		ArtistStatement: r.ArtistStatement,
		Email:           r.Email,
	}
}
