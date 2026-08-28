package restapi

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/maksimovyuriy/artfolio/backend/internal/controller/restapi/apierror"
	"github.com/maksimovyuriy/artfolio/backend/internal/controller/restapi/middleware"
	"github.com/maksimovyuriy/artfolio/backend/internal/controller/restapi/v1"
)

func NewRouter(
	controller *v1.Controller,
	authMiddleware *middleware.Auth,
	maxUploadBodySize int64,
	log *slog.Logger,
) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestLogger(log))
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

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	r.Mount("/v1", v1.NewRouter(controller, authMiddleware, maxUploadBodySize))

	return r
}
