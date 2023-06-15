package middleware

import (
	"log"
	"net/http"

	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"raspstore.github.io/users-service/internal/service"
)

type AuthorizationMiddleware interface {
	Apply(requiredPermissions ...string) func(http.Handler) http.Handler
}

type authorizationMiddleware struct {
	userService service.UserService
}

func NewAuthorizationMiddleware(userService service.UserService) AuthorizationMiddleware {
	return &authorizationMiddleware{userService: userService}
}

func (am *authorizationMiddleware) Apply(requiredPermissions ...string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			requesterId := r.Context().Value(UserIdKey).(string)

			requester, err := am.userService.GetUserById(requesterId)

			if err != nil {
				traceId := r.Context().Value(chiMiddleware.RequestIDKey).(string)
				log.Printf("[ERROR] - [%s]: Error while retrieving user information for id=%s: %s", traceId, requesterId, err.Error())
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
				return
			}

			for rp := range requiredPermissions {
				for p := range requester.Permissions {
					if rp == p {
						h.ServeHTTP(w, r)
					}
				}
			}
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		})
	}
}
