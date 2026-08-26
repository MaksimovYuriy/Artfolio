package v1

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/maksimovyuriy/artfolio/backend/internal/controller/restapi/jsonutil"
	"github.com/maksimovyuriy/artfolio/backend/internal/controller/restapi/v1/request"
	"github.com/maksimovyuriy/artfolio/backend/internal/controller/restapi/v1/response"
	"github.com/maksimovyuriy/artfolio/backend/internal/usecase"
)

func (c *Controller) getArtistProfile(w http.ResponseWriter, r *http.Request) {
	profile, err := c.artistProfile.Get(r.Context())
	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	_ = jsonutil.Encode(w, http.StatusOK, response.ArtistProfileFromEntity(profile))
}

func (c *Controller) updateArtistProfile(w http.ResponseWriter, r *http.Request) {
	var body request.UpdateArtistProfile
	if err := jsonutil.Decode(w, r, &body); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if err := c.artistProfile.Update(r.Context(), body.ArtistProfile()); err != nil {
		if errors.Is(err, usecase.ErrInvalidArtistProfile) {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
