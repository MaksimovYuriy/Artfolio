package v1

import (
	"net/http"

	"github.com/maksimovyuriy/artfolio/backend/internal/controller/restapi/apierror"
	"github.com/maksimovyuriy/artfolio/backend/internal/controller/restapi/jsonutil"
	"github.com/maksimovyuriy/artfolio/backend/internal/controller/restapi/v1/request"
	"github.com/maksimovyuriy/artfolio/backend/internal/controller/restapi/v1/response"
)

func (c *Controller) listSocialLinks(w http.ResponseWriter, r *http.Request) {
	links, err := c.socialLink.List(r.Context())
	if err != nil {
		apierror.Write(w, r, err)
		return
	}
	_ = jsonutil.Encode(w, http.StatusOK, response.AdminSocialLinksFromEntities(links))
}

func (c *Controller) replaceSocialLinks(w http.ResponseWriter, r *http.Request) {
	var body request.ReplaceSocialLinks
	if err := jsonutil.Decode(w, r, &body); err != nil {
		apierror.Write(w, r, apierror.InvalidRequest(err))
		return
	}
	if err := c.socialLink.Replace(r.Context(), body.Entities()); err != nil {
		apierror.Write(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
