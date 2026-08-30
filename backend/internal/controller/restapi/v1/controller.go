package v1

import (
	"github.com/maksimovyuriy/artfolio/backend/internal/controller/restapi/v1/response"
	"github.com/maksimovyuriy/artfolio/backend/internal/usecase"
)

type Controller struct {
	session       usecase.SessionUseCase
	artistProfile usecase.ArtistProfileUseCase
	artwork       usecase.ArtworkUseCase
	socialLink    usecase.SocialLinkUseCase
	contact       usecase.ContactUseCase
	artworkMapper response.ArtworkMapper
}

func NewController(
	session usecase.SessionUseCase,
	artistProfile usecase.ArtistProfileUseCase,
	artwork usecase.ArtworkUseCase,
	socialLink usecase.SocialLinkUseCase,
	contact usecase.ContactUseCase,
	artworkMapper response.ArtworkMapper,
) *Controller {
	return &Controller{
		session:       session,
		artistProfile: artistProfile,
		artwork:       artwork,
		socialLink:    socialLink,
		contact:       contact,
		artworkMapper: artworkMapper,
	}
}
