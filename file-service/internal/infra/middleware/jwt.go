package middleware

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/murilo-bracero/raspstore/file-service/internal/infra/config"
	"github.com/murilo-bracero/raspstore/file-service/internal/infra/response"
)

type userClaimsKeyType int

const UserClaimsCtxKey userClaimsKeyType = 101

const (
	tokenPrefix         = "Bearer"
	authorizationHeader = "Authorization"
)

var (
	ErrInvalidToken = errors.New("token is missing or is invalid")
)

func JWTMiddleware(config *config.Config) func(h http.Handler) http.Handler {
	ar := jwk.NewAutoRefresh(context.Background())

	ar.Configure(config.Auth.PublicKeyUrl, jwk.WithMinRefreshInterval(15*time.Minute))

	if _, err := ar.Refresh(context.Background(), config.Auth.PublicKeyUrl); err != nil {
		slog.Error("Failed to refresh JWT Tokens", "err", err)
	}

	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tkn, err := verifyJwt(r, ar, config.Auth.PublicKeyUrl)

			if err != nil {
				response.Unauthorized(w)
				return
			}

			ctx := context.WithValue(r.Context(), UserClaimsCtxKey, tkn)
			r = r.WithContext(ctx)

			h.ServeHTTP(w, r)
		})
	}
}

func verifyJwt(r *http.Request, ar *jwk.AutoRefresh, jwkUri string) (jwt.Token, error) {
	token, err := getTokenHeader(r)

	if err != nil {
		return nil, err
	}

	keyset, err := ar.Fetch(r.Context(), jwkUri)
	if err != nil {
		return nil, err
	}

	return jwt.Parse([]byte(token), jwt.WithKeySet(keyset), jwt.WithValidate(true))
}

func getTokenHeader(r *http.Request) (token string, err error) {
	header := r.Header.Get(authorizationHeader)

	split := strings.Split(header, " ")

	if len(split) != 2 {
		return "", ErrInvalidToken
	}

	prefix := split[0]

	if prefix != tokenPrefix {
		return "", ErrInvalidToken
	}

	return split[1], nil
}
