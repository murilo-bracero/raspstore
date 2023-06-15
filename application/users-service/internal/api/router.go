package api

import (
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"raspstore.github.io/users-service/internal/api/handler"
	"raspstore.github.io/users-service/internal/api/middleware"
)

const serviceBaseRoute = "/users-service/api/v1"
const userBaseRoute = serviceBaseRoute + "/users"
const userConfigRoute = serviceBaseRoute + "/config"
const adminBaseRoute = serviceBaseRoute + "/admin/users"

type Routes interface {
	MountRoutes() *chi.Mux
}

type routes struct {
	userHandler             handler.UserHandler
	userConfigHandler       handler.UserConfigHandler
	adminUserHandler        handler.AdminUserHandler
	authorizationMiddleware middleware.AuthorizationMiddleware
}

func NewRoutes(uc handler.UserHandler, uch handler.UserConfigHandler, auh handler.AdminUserHandler, am middleware.AuthorizationMiddleware) Routes {
	return &routes{userHandler: uc, userConfigHandler: uch, adminUserHandler: auh, authorizationMiddleware: am}
}

func (rt *routes) MountRoutes() *chi.Mux {

	router := chi.NewRouter()

	router.Use(middleware.Cors)
	router.Use(chiMiddleware.RequestID)
	router.Use(chiMiddleware.Logger)

	router.Route(userBaseRoute, func(r chi.Router) {
		r.With(middleware.Authentication).Get("/", rt.userHandler.ListUser)
		r.With(middleware.Authentication).Get("/{id}", rt.userHandler.GetUser)
		r.With(middleware.Authentication).Put("/{id}", rt.userHandler.UpdateUser)
		r.With(middleware.Authentication).Delete("/{id}", rt.userHandler.DeleteUser)
		r.Post("/", rt.userHandler.CreateUser)
	})

	router.Route(userConfigRoute, func(r chi.Router) {
		r.With(middleware.Authentication).Get("/", rt.userConfigHandler.GetUserConfigs)
		r.With(middleware.Authentication).Patch("/", rt.userConfigHandler.UpdateUserConfigs)
	})

	router.Route(adminBaseRoute, func(r chi.Router) {
		r.With(middleware.Authentication).With(rt.authorizationMiddleware.Apply("admin", "admin/user-create")).Post("/", rt.adminUserHandler.CreateUser)
	})

	return router

}
