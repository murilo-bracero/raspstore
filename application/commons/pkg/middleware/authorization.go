package middleware

import (
	"context"
	"log"
	"net/http"

	"raspstore.github.io/commons/pkg/service"
)

func Authorization(svc service.AuthService) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			token := r.Header.Get("Authorization")

			if token == "" {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			if uid, err := svc.Authenticate(token); err != nil {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			} else {
				ctx := context.WithValue(r.Context(), UserIdKey, uid.Uid)
				ctx = context.WithValue(ctx, UserRolesKey, uid.Roles)
				r = r.WithContext(ctx)
				log.Printf("[INFO] User %s is accessing resource %s", uid, r.RequestURI)
			}

			h.ServeHTTP(w, r)
		})
	}
}
