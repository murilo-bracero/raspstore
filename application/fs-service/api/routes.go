package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const serviceBaseRoute = "/fs-service"
const fileBaseRoute = serviceBaseRoute + "/files"

type Routes interface {
	MountRoutes() *chi.Mux
}

type routes struct {
	fc interface{}
}

func NewRoutes() Routes {
	return &routes{fc: nil}
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
