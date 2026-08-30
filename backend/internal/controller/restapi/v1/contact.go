package v1

import (
	"net/http"

	"github.com/maksimovyuriy/artfolio/backend/internal/controller/restapi/apierror"
	"github.com/maksimovyuriy/artfolio/backend/internal/controller/restapi/jsonutil"
	"github.com/maksimovyuriy/artfolio/backend/internal/controller/restapi/v1/request"
)

func (c *Controller) sendContactMessage(w http.ResponseWriter, r *http.Request) {
	var body request.ContactMessage
	if err := jsonutil.Decode(w, r, &body); err != nil {
		apierror.Write(w, r, apierror.InvalidRequest(err))
		return
	}
	if err := c.contact.Send(r.Context(), body.Entity()); err != nil {
		apierror.Write(w, r, err)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}
