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
	userHandler       handler.UserHandler
	userConfigHandler handler.UserConfigHandler
	adminUserHandler  handler.AdminUserHandler
}

func NewRoutes(uc handler.UserHandler, uch handler.UserConfigHandler, auh handler.AdminUserHandler) Routes {
	return &routes{userHandler: uc, userConfigHandler: uch, adminUserHandler: auh}
}

func (rt *routes) MountRoutes() *chi.Mux {

	router := chi.NewRouter()

	router.Use(middleware.Cors)
	router.Use(chiMiddleware.RequestID)
	router.Use(chiMiddleware.Logger)

	router.Route(userBaseRoute, func(r chi.Router) {
		r.With(middleware.Authorization).Get("/", rt.userHandler.ListUser)
		r.With(middleware.Authorization).Get("/{id}", rt.userHandler.GetUser)
		r.With(middleware.Authorization).Put("/{id}", rt.userHandler.UpdateUser)
		r.With(middleware.Authorization).Delete("/{id}", rt.userHandler.DeleteUser)
		r.Post("/", rt.userHandler.CreateUser)
	})

	router.Route(userConfigRoute, func(r chi.Router) {
		r.With(middleware.Authorization).Get("/", rt.userConfigHandler.GetUserConfigs)
		r.With(middleware.Authorization).Patch("/", rt.userConfigHandler.UpdateUserConfigs)
	})

	router.Route(adminBaseRoute, func(r chi.Router) {
		r.With(middleware.Authorization).With(middleware.Authentication("admin", "admin/user-create")).Post("/", rt.adminUserHandler.CreateUser)
	})

	return router
}
