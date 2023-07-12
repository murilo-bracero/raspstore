package api

import (
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	rMiddleware "github.com/murilo-bracero/raspstore/commons/pkg/middleware"
	"github.com/murilo-bracero/raspstore/commons/pkg/service"
	"github.com/murilo-bracero/raspstore/users-service/internal/api/handler"
	"github.com/murilo-bracero/raspstore/users-service/internal/api/middleware"
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
	authService       service.AuthService
}

func NewRoutes(uc handler.UserHandler, uch handler.UserConfigHandler, auh handler.AdminUserHandler, as service.AuthService) Routes {
	return &routes{userHandler: uc, userConfigHandler: uch, adminUserHandler: auh, authService: as}
}

func (rt *routes) MountRoutes() *chi.Mux {

	router := chi.NewRouter()

	router.Use(middleware.Cors)
	router.Use(chiMiddleware.RequestID)
	router.Use(chiMiddleware.Logger)

	authorizationMiddleware := rMiddleware.Authorization(rt.authService)

	router.Route(userBaseRoute, func(r chi.Router) {
		r.With(authorizationMiddleware).Get("/", rt.userHandler.ListUser)
		r.With(authorizationMiddleware).Get("/{id}", rt.userHandler.GetUser)
		r.With(authorizationMiddleware).Put("/{id}", rt.userHandler.UpdateUser)
		r.With(authorizationMiddleware).Delete("/{id}", rt.userHandler.DeleteUser)
		r.Post("/", rt.userHandler.CreateUser)
	})

	router.Route(userConfigRoute, func(r chi.Router) {
		r.With(authorizationMiddleware).Get("/", rt.userConfigHandler.GetUserConfigs)
		r.With(authorizationMiddleware).Patch("/", rt.userConfigHandler.UpdateUserConfigs)
	})

	router.Route(adminBaseRoute, func(r chi.Router) {
		r.With(authorizationMiddleware).With(rMiddleware.Authentication("admin", "admin/user-create")).Post("/", rt.adminUserHandler.CreateUser)
	})

	return router
}
