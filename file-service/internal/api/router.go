package api

import (
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	rmd "github.com/murilo-bracero/raspstore/commons/pkg/security/middleware"
	"github.com/murilo-bracero/raspstore/file-service/internal"
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
}

func NewFilesRouter(filesHandler handler.FilesHandler, uploadHandler handler.UploadHandler, downloadHandler handler.DownloadHandler) FilesRouter {
	return &filesRouter{filesHandler: filesHandler, uploadHandler: uploadHandler, downloadHandler: downloadHandler}
}

func (fr *filesRouter) MountRoutes() *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.Cors)
	router.Use(chiMiddleware.RequestID)
	router.Use(chiMiddleware.Logger)

	authMiddleware := rmd.Authorization(internal.PublicKey())

	router.Route(fileBaseRoute, func(r chi.Router) {
		r.With(authMiddleware).Get("/", fr.filesHandler.ListFiles)
		r.With(authMiddleware).Put("/{id}", fr.filesHandler.Update)
		r.With(authMiddleware).Delete("/{id}", fr.filesHandler.Delete)
	})

	router.With(authMiddleware).Post(uploadRoute, fr.uploadHandler.Upload)
	router.With(authMiddleware).Get(downloadRoute, fr.downloadHandler.Download)

	return router
}
