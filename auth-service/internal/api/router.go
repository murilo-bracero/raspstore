package api

import (
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/murilo-bracero/raspstore/auth-service/internal/api/handler"
	"github.com/murilo-bracero/raspstore/auth-service/internal/api/middleware"
	"github.com/murilo-bracero/raspstore/auth-service/internal/infra"
)

const baseRoute = "/idp"
const loginRoute = baseRoute + "/v1/login"
const profileRoute = baseRoute + "/v1/profile"
const adminRoute = baseRoute + "/v1/admin/users"

type CredentialsRouter interface {
	MountRoutes() *chi.Mux
}

type credentialsRouter struct {
	config         *infra.Config
	loginHandler   handler.LoginHandler
	profileHandler handler.ProfileHandler
	adminHandler   handler.AdminHandler
}

func NewRoutes(lh handler.LoginHandler, ph handler.ProfileHandler, ah handler.AdminHandler, config *infra.Config) CredentialsRouter {
	return &credentialsRouter{loginHandler: lh, profileHandler: ph, adminHandler: ah, config: config}
}

func (cr *credentialsRouter) MountRoutes() *chi.Mux {

	router := chi.NewRouter()

	router.Use(middleware.Cors)
	router.Use(chiMiddleware.RequestID)
	router.Use(chiMiddleware.Logger)

	router.Route(loginRoute, func(r chi.Router) {
		r.Post("/", cr.loginHandler.Login)
	})

	authorization := middleware.Authorization(cr.config)

	router.Route(profileRoute, func(r chi.Router) {
		r.Use(authorization)
		r.Get("/", cr.profileHandler.GetProfile)
		r.Put("/", cr.profileHandler.UpdateProfile)
		r.Delete("/", cr.profileHandler.DeleteProfile)
	})

	router.Route(adminRoute, func(r chi.Router) {
		//TODO:applu authentication
		r.Use(authorization)
		r.Post("/", cr.adminHandler.CreateUser)
		r.Get("/", cr.adminHandler.ListUsers)
		r.Get("/{userId}", cr.adminHandler.GetUserById)
		r.Put("/{userId}", cr.adminHandler.UpdateUserById)
		r.Delete("/{userId}", cr.adminHandler.DeleteUser)
	})

	return router
}
