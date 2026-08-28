package restapi

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/maksimovyuriy/artfolio/backend/internal/controller/restapi/middleware"
	"github.com/maksimovyuriy/artfolio/backend/internal/controller/restapi/v1"
	"github.com/maksimovyuriy/artfolio/backend/internal/controller/restapi/v1/response"
	"github.com/maksimovyuriy/artfolio/backend/internal/entity"
)

func TestRouterVersioningAndHealth(t *testing.T) {
	artistProfile := &routerArtistProfileUseCase{}
	session := &routerSessionUseCase{}
	controller := v1.NewController(session, artistProfile, nil, nil, response.NewArtworkMapper("/media"))
	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	router := NewRouter(controller, middleware.NewAuth(session), 1<<20, log)

	tests := []struct {
		path string
		want int
	}{
		{path: "/health", want: http.StatusOK},
		{path: "/v1/artist_profile", want: http.StatusOK},
		{path: "/artist_profile", want: http.StatusNotFound},
	}

	for _, test := range tests {
		response := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, test.path, nil)
		router.ServeHTTP(response, request)
		if response.Code != test.want {
			t.Errorf("GET %s status = %d, want %d", test.path, response.Code, test.want)
		}
	}
}

type routerArtistProfileUseCase struct{}

func (*routerArtistProfileUseCase) Get(context.Context) (entity.ArtistProfile, error) {
	return entity.ArtistProfile{Name: "Artist"}, nil
}

func (*routerArtistProfileUseCase) Update(context.Context, entity.ArtistProfile) error {
	return nil
}

type routerSessionUseCase struct{}

func (*routerSessionUseCase) Create(context.Context, string) (entity.Session, error) {
	return entity.Session{}, nil
}

func (*routerSessionUseCase) Authenticate(context.Context, string) (entity.AuthenticatedSession, error) {
	return entity.AuthenticatedSession{ID: 1, ActorID: 1}, nil
}

func (*routerSessionUseCase) Revoke(context.Context, string) (entity.AuthenticatedSession, error) {
	return entity.AuthenticatedSession{ID: 1, ActorID: 1}, nil
}
