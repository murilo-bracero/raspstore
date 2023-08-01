package api

import (
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	rMiddleware "github.com/murilo-bracero/raspstore/commons/pkg/middleware"
	"github.com/murilo-bracero/raspstore/commons/pkg/service"
	"github.com/murilo-bracero/raspstore/file-service/internal/api/handler"
	"github.com/murilo-bracero/raspstore/file-service/internal/api/middleware"
)

const serviceBaseRoute = "/file-service"
const fileBaseRoute = serviceBaseRoute + "/v1/files"
const uploadRoute = serviceBaseRoute + "/v1/uploads"
const downloadRoute = serviceBaseRoute + "/v1/downloads/{fileId}"

type FilesRouter interface {
	MountRoutes() *chi.Mux
}

type filesRouter struct {
	filesHandler    handler.FilesHandler
	uploadHandler   handler.UploadHandler
	downloadHandler handler.DownloadHandler
	authService     service.AuthService
}

func NewFilesRouter(filesHandler handler.FilesHandler, uploadHandler handler.UploadHandler, downloadHandler handler.DownloadHandler, as service.AuthService) FilesRouter {
	return &filesRouter{filesHandler: filesHandler, authService: as, uploadHandler: uploadHandler, downloadHandler: downloadHandler}
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

	router.With(authorizationMiddleware).Post(uploadRoute, fr.uploadHandler.Upload)
	router.With(authorizationMiddleware).Get(downloadRoute, fr.downloadHandler.Download)

	return router
}
