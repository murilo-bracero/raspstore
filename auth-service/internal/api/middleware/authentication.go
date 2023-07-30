package middleware

import (
	"net/http"

	"github.com/murilo-bracero/raspstore/auth-service/internal/model"
)

func Authentication(requiredPermissions ...string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims := r.Context().Value(UserClaimsCtxKey).(*model.UserClaims)

			if len(claims.Roles) == 0 {
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
				return
			}

			for _, rp := range requiredPermissions {
				for _, p := range claims.Roles {
					if rp == p {
						h.ServeHTTP(w, r)
						return
					}
				}
			}
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		})
	}
}
