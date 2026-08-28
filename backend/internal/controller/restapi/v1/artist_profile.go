package v1

import (
	"net/http"

	"github.com/maksimovyuriy/artfolio/backend/internal/controller/restapi/apierror"
	"github.com/maksimovyuriy/artfolio/backend/internal/controller/restapi/jsonutil"
	"github.com/maksimovyuriy/artfolio/backend/internal/controller/restapi/v1/request"
	"github.com/maksimovyuriy/artfolio/backend/internal/controller/restapi/v1/response"
)

func (c *Controller) getArtistProfile(w http.ResponseWriter, r *http.Request) {
	profile, err := c.artistProfile.Get(r.Context())
	if err != nil {
		apierror.Write(w, r, err)
		return
	}

	_ = jsonutil.Encode(w, http.StatusOK, response.ArtistProfileFromEntity(profile))
}

func (c *Controller) updateArtistProfile(w http.ResponseWriter, r *http.Request) {
	var body request.UpdateArtistProfile
	if err := jsonutil.Decode(w, r, &body); err != nil {
		apierror.Write(w, r, apierror.InvalidRequest(err))
		return
	}

	if err := c.artistProfile.Update(r.Context(), body.ArtistProfile()); err != nil {
		apierror.Write(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
