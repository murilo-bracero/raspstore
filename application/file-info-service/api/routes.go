package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"raspstore.github.io/file-manager/api/controller"
)

const serviceBaseRoute = "/file-info-service"
const fileBaseRoute = serviceBaseRoute + "/files"

type Routes interface {
	MountRoutes() *chi.Mux
}

type routes struct {
	fc controller.FilesController
}

func NewRoutes(fc controller.FilesController) Routes {
	return &routes{fc: fc}
}

func (rt *routes) MountRoutes() *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)

	router.Route(fileBaseRoute, func(r chi.Router) {
		r.Get("/", rt.fc.ListFiles)
		r.Put("/{id}", rt.fc.Update)
		r.Delete("/{id}", rt.fc.Delete)
	})

	return router

}
