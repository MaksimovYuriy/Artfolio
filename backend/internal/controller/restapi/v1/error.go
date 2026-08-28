package v1

import (
	"net/http"

	"github.com/maksimovyuriy/artfolio/backend/internal/controller/restapi/middleware"
)

func writeInternalError(w http.ResponseWriter, r *http.Request, err error) {
	middleware.RecordError(r.Context(), err)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
