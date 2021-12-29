package middleware

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	api "raspstore.github.io/authentication/api/dto"
	"raspstore.github.io/authentication/token"
)

type AuthMiddleware interface {
	Apply(h http.Handler) http.Handler
}

type authMiddleware struct {
	tm token.TokenManager
}

func NewAuthMiddleware(tm token.TokenManager) AuthMiddleware {
	return &authMiddleware{tm: tm}
}

const whitelistRoutes = "POST:/users,POST:/auth"

func (a *authMiddleware) Apply(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !isRouteWhitelisted(r.Method, r.RequestURI) {
			token := r.Header.Get("Authorization")

			if token == "" {
				w.WriteHeader(http.StatusUnauthorized)
				er := new(api.ErrorResponse)
				er.Message = "authorization header is missing"
				er.Code = "AM01"
				send(w, er)
				return
			}

			uid, err := a.tm.Verify(token)

			if err != nil {
				log.Println("token error : ", err.Error())
				w.WriteHeader(http.StatusUnauthorized)
				er := new(api.ErrorResponse)
				er.Message = "authorization token is denied"
				er.Code = "AM02"
				send(w, er)
				return
			}

			log.Printf("user %s is accessing resource %s", uid, r.RequestURI)
		}
		h.ServeHTTP(w, r)
	})
}

func isRouteWhitelisted(method string, path string) bool {
	routes := strings.Split(whitelistRoutes, ",")

	for _, route := range routes {
		aux := strings.Split(route, ":")
		m := aux[0]
		p := aux[1]

		if m == method && p == path {
			return true
		}
	}

	return false
}

func send(w http.ResponseWriter, obj interface{}) {
	w.Header().Set("Content-Type", "application/json")
	jsonResponse, err := json.Marshal(obj)
	if err != nil {
		return
	}
	w.Write(jsonResponse)
}
