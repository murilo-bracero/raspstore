package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/murilo-bracero/raspstore/auth-service/internal"
	"github.com/murilo-bracero/raspstore/auth-service/internal/infra"
	"github.com/murilo-bracero/raspstore/auth-service/internal/model"
	"github.com/murilo-bracero/raspstore/auth-service/internal/token"
)

type userClaimsKeyType int

const UserClaimsCtxKey userClaimsKeyType = 101

func Authorization(config *infra.Config) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			claims, err := getClaimsForRequest(config, r)

			if err != nil {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), UserClaimsCtxKey, claims)
			r = r.WithContext(ctx)

			h.ServeHTTP(w, r)
		})
	}
}

func getClaimsForRequest(config *infra.Config, r *http.Request) (*model.UserClaims, error) {
	claims, err := checkTokenInCookie(config, r)

	if err == nil {
		return claims, nil
	}

	return checkTokenInHeader(config, r)
}

func checkTokenInCookie(config *infra.Config, r *http.Request) (*model.UserClaims, error) {
	accessCookie, err := r.Cookie("access_token")

	if err != nil {
		return nil, err
	}

	accessToken := strings.ReplaceAll(accessCookie.Value, "Bearer ", "")

	return token.Verify(config, accessToken)
}

func checkTokenInHeader(config *infra.Config, r *http.Request) (*model.UserClaims, error) {
	header := r.Header.Get("Authorization")

	if header == "" {
		return nil, internal.ErrEmptyToken
	}

	accessToken := strings.ReplaceAll(header, "Bearer ", "")

	return token.Verify(config, accessToken)
}
