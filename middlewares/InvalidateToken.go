package middlewares

import (
	"net/http"
	"strings"

	"github.com/whoiswentz/goauth/infrastructure/cache"
)

func InvalidateToken(c *cache.Cache) Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {

			authorization := r.Header.Get("Authorization")
			if authorization == "" {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusBadRequest)
				return
			}

			token := strings.Split(authorization, " ")[1]

			c.Store(token, token)

			f(w, r)
		}
	}
}
