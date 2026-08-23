package restapi

import (
	"net/http"

	"github.com/go-chi/chi"
)

func NewRouter(
	sessionController *SessionController,
	artistProfileController *ArtistProfileController,
) http.Handler {
	r := chi.NewRouter()

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	r.Post("/admin/session", sessionController.Create)
	r.Get("/admin/session", sessionController.Verify)
	r.Delete("/admin/session", sessionController.Revoke)

	r.Get("/artist_profile", artistProfileController.Get)
	r.Put("/admin/artist_profile", artistProfileController.Update)

	return r
}
