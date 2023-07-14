package api

import (
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/murilo-bracero/raspstore/auth-service/internal/api/handler"
	"github.com/murilo-bracero/raspstore/auth-service/internal/api/middleware"
)

const baseRoute = "/auth-service"
const loginRoute = baseRoute + "/login"

type CredentialsRouter interface {
	MountRoutes() *chi.Mux
}

type credentialsRouter struct {
	credentialsHandler handler.CredentialsHandler
}

func NewRoutes(h handler.CredentialsHandler) CredentialsRouter {
	return &credentialsRouter{credentialsHandler: h}
}

func (cr *credentialsRouter) MountRoutes() *chi.Mux {

	router := chi.NewRouter()

	router.Use(middleware.Cors)
	router.Use(chiMiddleware.RequestID)
	router.Use(chiMiddleware.Logger)

	router.Route(loginRoute, func(r chi.Router) {
		r.Post("/", cr.credentialsHandler.Login)
	})

	return router
}
