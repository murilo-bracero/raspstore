package api

import (
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	rMiddleware "github.com/murilo-bracero/raspstore/commons/pkg/middleware"
	"github.com/murilo-bracero/raspstore/commons/pkg/service"
	"raspstore.github.io/fs-service/internal/api/handler"
	"raspstore.github.io/fs-service/internal/api/middleware"
)

const serviceBaseRoute = "/fs-service"
const fileBaseRoute = serviceBaseRoute + "/files"

type Routes interface {
	MountRoutes() *chi.Mux
}

type routes struct {
	fileServeHandler handler.FileServeHandler
	authService      service.AuthService
}

func NewRoutes(fsh handler.FileServeHandler, as service.AuthService) Routes {
	return &routes{fileServeHandler: fsh, authService: as}
}

func (rt *routes) MountRoutes() *chi.Mux {
	router := chi.NewRouter()

	authorizationMiddleware := rMiddleware.Authorization(rt.authService)

	router.Use(middleware.Cors)
	router.Use(chiMiddleware.RequestID)
	router.Use(chiMiddleware.Logger)
	router.Use(authorizationMiddleware)

	router.Route(fileBaseRoute, func(r chi.Router) {
		r.Get("/{fileId}", rt.fileServeHandler.Download)
		r.Post("/", rt.fileServeHandler.Upload)
	})

	return router

}
