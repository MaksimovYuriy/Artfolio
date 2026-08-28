package v1

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/maksimovyuriy/artfolio/backend/internal/controller/restapi/jsonutil"
	"github.com/maksimovyuriy/artfolio/backend/internal/controller/restapi/middleware"
	"github.com/maksimovyuriy/artfolio/backend/internal/controller/restapi/v1/request"
	"github.com/maksimovyuriy/artfolio/backend/internal/usecase"
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
	middleware.RecordAuthentication(r.Context(), adminSession.ActorID, adminSession.ID)

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

	adminSession, err := c.session.Authenticate(r.Context(), cookie.Value)
	if errors.Is(err, usecase.ErrInvalidSession) {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	if err != nil {
		writeInternalError(w, r, err)
		return
	}
	middleware.RecordAuthentication(r.Context(), adminSession.ActorID, adminSession.ID)

	w.WriteHeader(http.StatusNoContent)
}

func (c *Controller) revokeSession(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(middleware.SessionCookieName)
	if err == nil {
		adminSession, err := c.session.Revoke(r.Context(), cookie.Value)
		if err != nil && !errors.Is(err, usecase.ErrInvalidSession) {
			writeInternalError(w, r, err)
			return
		}
		if err == nil {
			middleware.RecordAuthentication(r.Context(), adminSession.ActorID, adminSession.ID)
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
