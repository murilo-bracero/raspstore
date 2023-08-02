package api

import (
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	rmd "github.com/murilo-bracero/raspstore/commons/pkg/security/middleware"
	"github.com/murilo-bracero/raspstore/idp/internal/api/handler"
	"github.com/murilo-bracero/raspstore/idp/internal/api/middleware"
	"github.com/murilo-bracero/raspstore/idp/internal/infra"
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

	authorization := rmd.Authorization(cr.config.TokenPublicKey)

	router.Route(profileRoute, func(r chi.Router) {
		r.Use(authorization)
		r.Get("/", cr.profileHandler.GetProfile)
		r.Put("/", cr.profileHandler.UpdateProfile)
		r.Delete("/", cr.profileHandler.DeleteProfile)
	})

	router.Route(adminRoute, func(r chi.Router) {
		r.Use(authorization)
		r.With(rmd.Authentication("admin", "admin/create-user")).Post("/", cr.adminHandler.CreateUser)
		r.With(rmd.Authentication("admin", "admin/list-user", "admin/get-user")).Get("/", cr.adminHandler.ListUsers)
		r.With(rmd.Authentication("admin", "admin/list-user", "admin/get-user")).Get("/{userId}", cr.adminHandler.GetUserById)
		r.With(rmd.Authentication("admin", "admin/update-user")).Put("/{userId}", cr.adminHandler.UpdateUserById)
		r.With(rmd.Authentication("admin", "admin/delete-user")).Delete("/{userId}", cr.adminHandler.DeleteUser)
	})

	return router
}
