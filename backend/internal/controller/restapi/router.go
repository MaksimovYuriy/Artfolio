package restapi

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/maksimovyuriy/artfolio/backend/internal/controller/restapi/middleware"
	"github.com/maksimovyuriy/artfolio/backend/internal/controller/restapi/v1"
)

func NewRouter(
	controller *v1.Controller,
	authMiddleware *middleware.Auth,
	maxUploadBodySize int64,
) http.Handler {
	r := chi.NewRouter()

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	r.Mount("/v1", v1.NewRouter(controller, authMiddleware, maxUploadBodySize))

	return r
}
