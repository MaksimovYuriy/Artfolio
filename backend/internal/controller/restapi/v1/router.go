package v1

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/maksimovyuriy/artfolio/backend/internal/controller/restapi/middleware"
)

func NewRouter(controller *Controller, authMiddleware *middleware.Auth, maxUploadBodySize int64) http.Handler {
	r := chi.NewRouter()

	r.Get("/artist_profile", controller.getArtistProfile)
	r.Get("/artworks", controller.listPublishedArtworks)

	r.Route("/admin", func(r chi.Router) {
		r.Post("/session", controller.createSession)
		r.Get("/session", controller.verifySession)
		r.Delete("/session", controller.revokeSession)

		r.Group(func(r chi.Router) {
			r.Use(authMiddleware.VerifySession)

			r.Put("/artist_profile", controller.updateArtistProfile)

			r.Get("/artworks", controller.listAllArtworks)
			r.With(middleware.MaxBodySize(maxUploadBodySize)).Post("/artworks", controller.createArtwork)
			r.With(middleware.MaxBodySize(maxUploadBodySize)).Put("/artworks/:id", controller.updateArtwork)
			r.Delete("/artworks/:id", controller.deleteArtwork)
		})
	})

	return r
}
