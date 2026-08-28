package v1

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/maksimovyuriy/artfolio/backend/internal/controller/restapi/apierror"
	"github.com/maksimovyuriy/artfolio/backend/internal/controller/restapi/middleware"
)

func NewRouter(controller *Controller, authMiddleware *middleware.Auth, maxUploadBodySize int64) http.Handler {
	r := chi.NewRouter()
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		apierror.Write(w, r, apierror.WithStatus(http.StatusNotFound, errors.New("route not found")))
	})
	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		apierror.Write(w, r, apierror.New(
			http.StatusMethodNotAllowed,
			"method_not_allowed",
			"Method not allowed",
			errors.New("method not allowed"),
		))
	})

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
			r.Put("/artworks/order", controller.reorderArtworks)
			r.With(middleware.MaxBodySize(maxUploadBodySize)).Put("/artworks/{id}", controller.updateArtwork)
			r.Delete("/artworks/{id}", controller.deleteArtwork)

			r.Get("/social_links", controller.listSocialLinks)
			r.Put("/social_links", controller.replaceSocialLinks)
		})
	})

	return r
}
