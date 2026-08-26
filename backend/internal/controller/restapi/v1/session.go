package v1

import (
	"net/http"
	"strings"
	"time"

	"github.com/maksimovyuriy/artfolio/backend/internal/controller/restapi/jsonutil"
	"github.com/maksimovyuriy/artfolio/backend/internal/controller/restapi/middleware"
	"github.com/maksimovyuriy/artfolio/backend/internal/controller/restapi/v1/request"
)

func (c *Controller) createSession(w http.ResponseWriter, r *http.Request) {
	var body request.CreateSession
	if err := jsonutil.Decode(w, r, &body); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	body.AccessKey = strings.TrimSpace(body.AccessKey)
	if body.AccessKey == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	adminSession, err := c.session.Create(r.Context(), body.AccessKey)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     middleware.SessionCookieName,
		Value:    adminSession.Token,
		Path:     "/",
		Expires:  adminSession.ExpiresAt,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	w.WriteHeader(http.StatusNoContent)
}

func (c *Controller) verifySession(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(middleware.SessionCookieName)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	valid, err := c.session.Verify(r.Context(), cookie.Value)
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

func (c *Controller) revokeSession(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(middleware.SessionCookieName)
	if err == nil {
		if err := c.session.Revoke(r.Context(), cookie.Value); err != nil {
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
