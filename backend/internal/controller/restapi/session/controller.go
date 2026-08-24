package session

import (
	"net/http"
	"strings"
	"time"

	sessiondto "github.com/maksimovyuriy/artfolio/backend/internal/controller/restapi/dto/session"
	"github.com/maksimovyuriy/artfolio/backend/internal/controller/restapi/jsonutil"
	"github.com/maksimovyuriy/artfolio/backend/internal/controller/restapi/middleware"
	"github.com/maksimovyuriy/artfolio/backend/internal/usecase"
)

type Controller struct {
	useCase usecase.SessionUseCase
}

func New(useCase usecase.SessionUseCase) *Controller {
	return &Controller{useCase: useCase}
}

func (c *Controller) Create(w http.ResponseWriter, r *http.Request) {
	var request sessiondto.CreateSessionRequest
	if err := jsonutil.Decode(w, r, &request); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	request.AccessKey = strings.TrimSpace(request.AccessKey)
	if request.AccessKey == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	session, err := c.useCase.Create(r.Context(), request.AccessKey)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     middleware.SessionCookieName,
		Value:    session.Token,
		Path:     "/",
		Expires:  session.ExpiresAt,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	w.WriteHeader(http.StatusNoContent)
}

func (c *Controller) Verify(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(middleware.SessionCookieName)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	valid, err := c.useCase.Verify(r.Context(), cookie.Value)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if !valid {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (c *Controller) Revoke(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(middleware.SessionCookieName)
	if err == nil {
		if err := c.useCase.Revoke(r.Context(), cookie.Value); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}

	http.SetCookie(w, &http.Cookie{
		Name:     middleware.SessionCookieName,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	w.WriteHeader(http.StatusNoContent)
}
