package api

import (
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	rMiddleware "github.com/murilo-bracero/raspstore/commons/pkg/middleware"
	"github.com/murilo-bracero/raspstore/commons/pkg/service"
	"github.com/murilo-bracero/raspstore/file-info-service/internal/api/handler"
	"github.com/murilo-bracero/raspstore/file-info-service/internal/api/middleware"
)

const serviceBaseRoute = "/file-info-service"
const fileBaseRoute = serviceBaseRoute + "/files"

type FilesRouter interface {
	MountRoutes() *chi.Mux
}

type filesRouter struct {
	filesHandler handler.FilesHandler
	authService  service.AuthService
}

func NewFilesRouter(filesHandler handler.FilesHandler, as service.AuthService) FilesRouter {
	return &filesRouter{filesHandler: filesHandler, authService: as}
}

func (fr *filesRouter) MountRoutes() *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.Cors)
	router.Use(chiMiddleware.RequestID)
	router.Use(chiMiddleware.Logger)

	authorizationMiddleware := rMiddleware.Authorization(fr.authService)

	router.Route(fileBaseRoute, func(r chi.Router) {
		r.With(authorizationMiddleware).Get("/", fr.filesHandler.ListFiles)
		r.With(authorizationMiddleware).Put("/{id}", fr.filesHandler.Update)
		r.With(authorizationMiddleware).Delete("/{id}", fr.filesHandler.Delete)
	})

	return router
}
