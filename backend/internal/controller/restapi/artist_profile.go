package restapi

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/maksimovyuriy/artfolio/backend/internal/controller/restapi/dto"
	"github.com/maksimovyuriy/artfolio/backend/internal/usecase"
)

type ArtistProfileController struct {
	usecase usecase.ArtistProfileUseCase
}

func NewArtistProfileController(
	usecase usecase.ArtistProfileUseCase,
) *ArtistProfileController {
	return &ArtistProfileController{usecase: usecase}
}

func (c *ArtistProfileController) Get(w http.ResponseWriter, r *http.Request) {
	profile, err := c.usecase.Get(r.Context())
	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	response := dto.ArtistProfileResponseFromEntity(profile)
	_ = encodeJSON(w, http.StatusOK, response)
}

func (c *ArtistProfileController) Update(w http.ResponseWriter, r *http.Request) {
	var request dto.UpdateArtistProfileRequest
	if err := decodeJSON(w, r, &request); err != nil {
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
