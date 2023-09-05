package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
)

type userClaimsKeyType int

const UserClaimsCtxKey userClaimsKeyType = 101

const (
	tokenPrefix         = "Bearer"
	authorizationHeader = "Authorization"
	keycloakUrl         = "http://localhost:30101/realms/raspstore-dev"
)

var (
	ErrInvalidToken = errors.New("token is missing or is invalid")
)

func JWTMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tkn, err := verifyJwt(r)

		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserClaimsCtxKey, tkn)
		r = r.WithContext(ctx)

		h.ServeHTTP(w, r)
	})
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

func verifyJwt(r *http.Request) (jwt.Token, error) {
	token, err := getTokenHeader(r)

	if err != nil {
		return nil, err
	}

	keyset, err := jwk.Fetch(r.Context(), keycloakUrl+"/protocol/openid-connect/certs")
	if err != nil {
		return nil, err
	}

	return jwt.Parse([]byte(token), jwt.WithKeySet(keyset), jwt.WithValidate(true))
}
