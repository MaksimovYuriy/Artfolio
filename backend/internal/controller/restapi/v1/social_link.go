package v1

import (
	"errors"
	"net/http"

	"github.com/maksimovyuriy/artfolio/backend/internal/controller/restapi/jsonutil"
	"github.com/maksimovyuriy/artfolio/backend/internal/controller/restapi/v1/request"
	"github.com/maksimovyuriy/artfolio/backend/internal/controller/restapi/v1/response"
	"github.com/maksimovyuriy/artfolio/backend/internal/entity"
)

func (c *Controller) listSocialLinks(w http.ResponseWriter, r *http.Request) {
	links, err := c.socialLink.List(r.Context())
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	_ = jsonutil.Encode(w, http.StatusOK, response.AdminSocialLinksFromEntities(links))
}

func (c *Controller) replaceSocialLinks(w http.ResponseWriter, r *http.Request) {
	var body request.ReplaceSocialLinks
	if err := jsonutil.Decode(w, r, &body); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if err := c.socialLink.Replace(r.Context(), body.Entities()); err != nil {
		if errors.Is(err, entity.ErrValidation) {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
