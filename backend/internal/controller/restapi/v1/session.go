package v1

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/maksimovyuriy/artfolio/backend/internal/controller/restapi/apierror"
	"github.com/maksimovyuriy/artfolio/backend/internal/controller/restapi/jsonutil"
	"github.com/maksimovyuriy/artfolio/backend/internal/controller/restapi/middleware"
	"github.com/maksimovyuriy/artfolio/backend/internal/controller/restapi/v1/request"
	"github.com/maksimovyuriy/artfolio/backend/internal/usecase"
)

func (c *Controller) createSession(w http.ResponseWriter, r *http.Request) {
	var body request.CreateSession
	if err := jsonutil.Decode(w, r, &body); err != nil {
		apierror.Write(w, r, apierror.InvalidRequest(err))
		return
	}

	body.AccessKey = strings.TrimSpace(body.AccessKey)
	if body.AccessKey == "" {
		apierror.Write(w, r, apierror.InvalidRequest(errors.New("access key is required")))
		return
	}

	adminSession, err := c.session.Create(r.Context(), body.AccessKey)
	if err != nil {
		apierror.Write(w, r, err)
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
		apierror.Write(w, r, usecase.ErrInvalidSession)
		return
	}

	valid, err := c.session.Verify(r.Context(), cookie.Value)
	if err != nil {
		apierror.Write(w, r, err)
		return
	}
	if !valid {
		apierror.Write(w, r, usecase.ErrInvalidSession)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (c *Controller) revokeSession(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(middleware.SessionCookieName)
	if err == nil {
		if err := c.session.Revoke(r.Context(), cookie.Value); err != nil {
			apierror.Write(w, r, err)
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
