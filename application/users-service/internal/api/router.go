package api

import (
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"raspstore.github.io/users-service/internal/api/handler"
	"raspstore.github.io/users-service/internal/api/middleware"
)

const serviceBaseRoute = "/users-service"
const userBaseRoute = serviceBaseRoute + "/users"

type Routes interface {
	MountRoutes() *chi.Mux
}

type routes struct {
	userHandler handler.UserHandler
}

func NewRoutes(uc handler.UserHandler) Routes {
	return &routes{userHandler: uc}
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

	return router

}
