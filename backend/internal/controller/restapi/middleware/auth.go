package middleware

import (
	"errors"
	"net/http"

	"github.com/maksimovyuriy/artfolio/backend/internal/usecase"
)

const SessionCookieName = "artfolio_session"

type Auth struct {
	usecase usecase.SessionUseCase
}

func NewAuth(usecase usecase.SessionUseCase) *Auth {
	return &Auth{usecase: usecase}
}

func (a *Auth) VerifySession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(SessionCookieName)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		session, err := a.usecase.Authenticate(r.Context(), cookie.Value)
		if errors.Is(err, usecase.ErrInvalidSession) {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		if err != nil {
			RecordError(r.Context(), err)
			http.Error(
				w,
				http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError,
			)
			return
		}
		RecordAuthentication(r.Context(), session.ActorID, session.ID)

		next.ServeHTTP(w, r)
	})
}
