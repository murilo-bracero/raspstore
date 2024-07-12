package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/murilo-bracero/raspstore/file-service/internal/application/usecase"
	"github.com/murilo-bracero/raspstore/file-service/internal/infra/config"
	"github.com/murilo-bracero/raspstore/file-service/internal/infra/handler"
)

func StartApiServer(config *config.Config,
	luc usecase.ListFilesUseCase,
	uuc usecase.UpdateFileUseCase,
	duc usecase.DeleteFileUseCase,
	upc usecase.UploadFileUseCase,
	downloadUc usecase.DownloadFileUseCase,
	createUc usecase.CreateFileUseCase,
	getFileUc usecase.GetFileUseCase) {
	filesHandler := handler.NewFilesHandler(luc, uuc, duc)

	uploadHanler := handler.NewUploadHandler(config, upc, createUc)

	downloadHandler := handler.NewDownloadHandler(downloadUc, getFileUc)

	router := NewFilesRouter(config, filesHandler, uploadHanler, downloadHandler).MountRoutes()
	http.Handle("/", router)
	slog.Info("File Manager REST API runing", "port", config.Server.Port)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", config.Server.Port), router); err != nil {
		slog.Error("Could not start File Manager REST server", "error", err)
		os.Exit(1)
	}
}
