package middlewares

import (
	"net/http"
	"strings"

	"github.com/whoiswentz/goauth/helpers"
	"github.com/whoiswentz/goauth/infrastructure/cache"
)

func RequireToken(c *cache.Cache) Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {

			authorization := r.Header.Get("Authorization")
			if authorization == "" {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusBadRequest)
				return
			}

			token := strings.Split(authorization, " ")[1]

			_, load := c.Load(token)
			if load {
				http.Error(w, http.StatusText(http.StatusForbidden), http.StatusBadRequest)
				return
			}

			if _, err := helpers.ValidateToken(token); err != nil {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}

			f(w, r)
		}
	}
}
