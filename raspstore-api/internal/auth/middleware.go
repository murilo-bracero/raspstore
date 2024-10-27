package auth

import (
	"context"
	"net/http"

	"github.com/murilo-bracero/raspstore/file-service/internal/infra/validator"
)

const authorizationHeader = "Authorization"

type userClaimsKeyType int

const UserClaimsCtxKey userClaimsKeyType = 101

func TokenMiddleware(validator *validator.JWTValidator) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tkn, err := validator.Validate(r.Context(), r.Header.Get(authorizationHeader))

			if err != nil {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), UserClaimsCtxKey, *tkn)
			r = r.WithContext(ctx)

			h.ServeHTTP(w, r)
		})
	}
}
