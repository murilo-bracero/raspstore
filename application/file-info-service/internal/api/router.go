package api

import (
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"raspstore.github.io/file-manager/internal/api/handler"
	"raspstore.github.io/file-manager/internal/api/middleware"
)

const serviceBaseRoute = "/file-info-service"
const fileBaseRoute = serviceBaseRoute + "/files"

type FilesRouter interface {
	MountRoutes() *chi.Mux
}

type filesRouter struct {
	filesHandler handler.FilesHandler
}

func NewFilesRouter(filesHandler handler.FilesHandler) FilesRouter {
	return &filesRouter{filesHandler: filesHandler}
}

func (fr *filesRouter) MountRoutes() *chi.Mux {
	router := chi.NewRouter()

	router.Use(chiMiddleware.RequestID)
	router.Use(chiMiddleware.Logger)

	router.Route(fileBaseRoute, func(r chi.Router) {
		r.With(middleware.AuthenticationMiddleware).Get("/", fr.filesHandler.ListFiles)
		r.With(middleware.AuthenticationMiddleware).Put("/{id}", fr.filesHandler.Update)
		r.With(middleware.AuthenticationMiddleware).Delete("/{id}", fr.filesHandler.Delete)
	})

	return router
}
