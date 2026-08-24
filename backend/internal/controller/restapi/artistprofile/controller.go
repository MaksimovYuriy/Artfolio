package artistprofile

import (
	"database/sql"
	"errors"
	"net/http"

	artistprofiledto "github.com/maksimovyuriy/artfolio/backend/internal/controller/restapi/dto/artistprofile"
	"github.com/maksimovyuriy/artfolio/backend/internal/controller/restapi/jsonutil"
	"github.com/maksimovyuriy/artfolio/backend/internal/usecase"
)

type Controller struct {
	usecase usecase.ArtistProfileUseCase
}

func New(
	usecase usecase.ArtistProfileUseCase,
) *Controller {
	return &Controller{usecase: usecase}
}

func (c *Controller) Get(w http.ResponseWriter, r *http.Request) {
	profile, err := c.usecase.Get(r.Context())
	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	response := artistprofiledto.ArtistProfileResponseFromEntity(profile)
	_ = jsonutil.Encode(w, http.StatusOK, response)
}

func (c *Controller) Update(w http.ResponseWriter, r *http.Request) {
	var request artistprofiledto.UpdateArtistProfileRequest
	if err := jsonutil.Decode(w, r, &request); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if err := c.usecase.Update(r.Context(), request.ArtistProfile()); err != nil {
		if errors.Is(err, usecase.ErrInvalidArtistProfile) {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
