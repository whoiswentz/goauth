package middlewares

import (
	"net/http"
	"strings"

	"github.com/whoiswentz/goauth/helpers"
)

func RequireToken() Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {

			authorization := r.Header.Get("Authorization")
			if authorization == "" {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusBadRequest)
				return
			}

			token := strings.Split(authorization, " ")[0]
			if _, err := helpers.ValidateToken(token); err != nil {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}

			f(w, r)
		}
	}
}
