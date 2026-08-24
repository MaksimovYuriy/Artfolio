package restapi

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/maksimovyuriy/artfolio/backend/internal/controller/restapi/artistprofile"
	"github.com/maksimovyuriy/artfolio/backend/internal/controller/restapi/artwork"
	"github.com/maksimovyuriy/artfolio/backend/internal/controller/restapi/middleware"
	"github.com/maksimovyuriy/artfolio/backend/internal/controller/restapi/session"
)

func NewRouter(
	sessionController *session.Controller,
	artistProfileController *artistprofile.Controller,
	artworkController *artwork.Controller,
	authMiddleware *middleware.Auth,
) http.Handler {
	r := chi.NewRouter()

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	r.Get("/artist_profile", artistProfileController.Get)
	r.Get("/artworks", artworkController.ListPublished)

	r.Route("/admin", func(r chi.Router) {
		r.Post("/session", sessionController.Create)
		r.Get("/session", sessionController.Verify)
		r.Delete("/session", sessionController.Revoke)

		r.Group(func(r chi.Router) {
			r.Use(authMiddleware.VerifySession)

			r.Put("/artist_profile", artistProfileController.Update)

			r.Get("/artworks", artworkController.ListAll)
			r.Post("/artworks", artworkController.Create)
			r.Put("/artworks/:id", artworkController.Update)
			r.Delete("/artworks/:id", artworkController.Delete)
		})
	})

	return r
}
