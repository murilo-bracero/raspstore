package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"raspstore.github.io/fs-service/api/controller"
)

const serviceBaseRoute = "/fs-service"
const fileBaseRoute = serviceBaseRoute + "/files"

type Routes interface {
	MountRoutes() *chi.Mux
}

type routes struct {
	fc controller.FileServeController
}

func NewRoutes(fc controller.FileServeController) Routes {
	return &routes{fc: fc}
}

func (rt *routes) MountRoutes() *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)

	router.Route(fileBaseRoute, func(r chi.Router) {
		r.Get("/{id}", rt.fc.Download)
		r.Post("/", rt.fc.Upload)
	})

	return router

}
