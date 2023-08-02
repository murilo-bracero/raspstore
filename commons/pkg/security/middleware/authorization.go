package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"aidanwoods.dev/go-paseto"
	"github.com/murilo-bracero/raspstore/commons/pkg/object"
)

type userClaimsKeyType int

const UserClaimsCtxKey userClaimsKeyType = 101

const (
	audience            = "account"
	tokenPrefix         = "Bearer "
	authorizationHeader = "Authorization"
	accessTokenCookie   = "access_token"
)

var (
	ErrEmptyToken = errors.New("Authorization header is empty")
)

func Authorization(publicKey string) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			claims, err := getClaimsForRequest(publicKey, r)

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

func getClaimsForRequest(publicKey string, r *http.Request) (*object.Claims, error) {
	claims, err := checkTokenInCookie(publicKey, r)

	if err == nil {
		return claims, nil
	}

	return checkTokenInHeader(publicKey, r)
}

func checkTokenInCookie(publicKey string, r *http.Request) (*object.Claims, error) {
	accessCookie, err := r.Cookie(accessTokenCookie)

	if err != nil {
		return nil, err
	}

	accessToken := strings.ReplaceAll(accessCookie.Value, tokenPrefix, "")

	return verify(publicKey, accessToken)
}

func checkTokenInHeader(publicKey string, r *http.Request) (*object.Claims, error) {
	header := r.Header.Get(authorizationHeader)

	if header == "" {
		return nil, ErrEmptyToken
	}

	accessToken := strings.ReplaceAll(header, tokenPrefix, "")

	return verify(publicKey, accessToken)
}

func verify(publicKey string, token string) (*object.Claims, error) {
	parser := paseto.NewParser()

	parser.AddRule(paseto.ForAudience(audience))
	parser.AddRule(paseto.NotExpired())
	parser.AddRule(paseto.ValidAt(time.Now()))

	key, err := paseto.NewV4AsymmetricPublicKeyFromHex(publicKey)

	if err != nil {
		return nil, err
	}

	pToken, err := parser.ParseV4Public(key, token, nil)

	if err != nil {
		return nil, err
	}

	return object.NewClaimsFromToken(pToken), nil
}
