package middleware

import (
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

		valid, err := a.usecase.Verify(r.Context(), cookie.Value)
		if err != nil {
			http.Error(
				w,
				http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError,
			)
			return
		}
		if !valid {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
