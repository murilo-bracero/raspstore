package server

import (
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/murilo-bracero/raspstore/file-service/internal/infra/config"
	"github.com/murilo-bracero/raspstore/file-service/internal/infra/handler"
	"github.com/murilo-bracero/raspstore/file-service/internal/infra/middleware"
)

const serviceBaseRoute = "/file-service"
const fileBaseRoute = serviceBaseRoute + "/v1/files"
const uploadRoute = serviceBaseRoute + "/v1/uploads"
const downloadRoute = serviceBaseRoute + "/v1/downloads/{fileId}"

type FilesRouter interface {
	MountRoutes() *chi.Mux
}

type filesRouter struct {
	config          *config.Config
	filesHandler    handler.FilesHandler
	uploadHandler   handler.UploadHandler
	downloadHandler handler.DownloadHandler
}

func NewFilesRouter(config *config.Config, filesHandler handler.FilesHandler, uploadHandler handler.UploadHandler, downloadHandler handler.DownloadHandler) FilesRouter {
	return &filesRouter{config: config, filesHandler: filesHandler, uploadHandler: uploadHandler, downloadHandler: downloadHandler}
}

func (fr *filesRouter) MountRoutes() *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.Cors)
	router.Use(chiMiddleware.RequestID)
	router.Use(chiMiddleware.Logger)
	router.Use(middleware.JWTMiddleware(fr.config))

	router.Route(fileBaseRoute, func(r chi.Router) {
		r.Get("/", fr.filesHandler.ListFiles)
		r.Get("/{id}", fr.filesHandler.FindById)
		r.Put("/{id}", fr.filesHandler.Update)
		r.Delete("/{id}", fr.filesHandler.Delete)
	})

	router.Post(uploadRoute, fr.uploadHandler.Upload)
	router.Get(downloadRoute, fr.downloadHandler.Download)

	return router
}
