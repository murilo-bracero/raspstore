package api

import (
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/murilo-bracero/raspstore/auth-service/internal/api/handler"
	"github.com/murilo-bracero/raspstore/auth-service/internal/api/middleware"
)

const baseRoute = "/idp"
const loginRoute = baseRoute + "/v1/login"
const profileRoute = baseRoute + "/v1/profile"

type CredentialsRouter interface {
	MountRoutes() *chi.Mux
}

type credentialsRouter struct {
	loginHandler   handler.LoginHandler
	profileHandler handler.ProfileHandler
}

func NewRoutes(lh handler.LoginHandler, ph handler.ProfileHandler) CredentialsRouter {
	return &credentialsRouter{loginHandler: lh, profileHandler: ph}
}

func (cr *credentialsRouter) MountRoutes() *chi.Mux {

	router := chi.NewRouter()

	router.Use(middleware.Cors)
	router.Use(chiMiddleware.RequestID)
	router.Use(chiMiddleware.Logger)

	router.Route(loginRoute, func(r chi.Router) {
		r.Post("/", cr.loginHandler.Login)
	})

	router.Route(profileRoute, func(r chi.Router) {
		r.Get("/", cr.profileHandler.GetProfile)
		r.Put("/", cr.profileHandler.UpdateProfile)
		r.Delete("/", cr.profileHandler.DeleteProfile)
	})

	return router
}
