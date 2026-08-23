package restapi

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/maksimovyuriy/artfolio/backend/internal/usecase"
)

const sessionCookieName = "artfolio_session"

type SessionController struct {
	useCase usecase.SessionUseCase
}

func NewSessionController(useCase usecase.SessionUseCase) *SessionController {
	return &SessionController{useCase: useCase}
}

type createSessionRequest struct {
	AccessKey string `json:"accessKey"`
}

func (c *SessionController) Create(w http.ResponseWriter, r *http.Request) {
	var request createSessionRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
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
		Name:     sessionCookieName,
		Value:    session.Token,
		Path:     "/",
		Expires:  session.ExpiresAt,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	w.WriteHeader(http.StatusNoContent)
}

func (c *SessionController) Verify(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(sessionCookieName)
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

func (c *SessionController) Revoke(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(sessionCookieName)
	if err == nil {
		if err := c.useCase.Revoke(r.Context(), cookie.Value); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}

	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
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
