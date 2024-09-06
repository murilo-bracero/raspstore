package middleware

import (
	"context"
	"net/http"

	"github.com/murilo-bracero/raspstore/file-service/internal/infra/response"
	"github.com/murilo-bracero/raspstore/file-service/internal/infra/validator"
)

type userClaimsKeyType int

const UserClaimsCtxKey userClaimsKeyType = 101

const authorizationHeader = "Authorization"

func JWTMiddleware(validator *validator.JWTValidator) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tkn, err := validator.Validate(r.Context(), r.Header.Get(authorizationHeader))

			if err != nil {
				response.Unauthorized(w)
				return
			}

			ctx := context.WithValue(r.Context(), UserClaimsCtxKey, *tkn)
			r = r.WithContext(ctx)

			h.ServeHTTP(w, r)
		})
	}
}
