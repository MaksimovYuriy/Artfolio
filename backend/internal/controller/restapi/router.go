package restapi

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/maksimovyuriy/artfolio/backend/internal/controller/restapi/middleware"
)

func NewRouter(
	sessionController *SessionController,
	artistProfileController *ArtistProfileController,
	authMiddleware *middleware.Auth,
) http.Handler {
	r := chi.NewRouter()

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	r.Get("/artist_profile", artistProfileController.Get)

	r.Route("/admin", func(r chi.Router) {
		r.Post("/session", sessionController.Create)
		r.Get("/session", sessionController.Verify)
		r.Delete("/session", sessionController.Revoke)

		r.Group(func(r chi.Router) {
			r.Use(authMiddleware.VerifySession)

			r.Put("/artist_profile", artistProfileController.Update)
		})
	})

	return r
}
