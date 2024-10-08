package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/murilo-bracero/raspstore/file-service/internal/application/facade"
	"github.com/murilo-bracero/raspstore/file-service/internal/application/usecase"
	"github.com/murilo-bracero/raspstore/file-service/internal/infra/config"
	"github.com/murilo-bracero/raspstore/file-service/internal/infra/handler"
)

func StartApiServer(config *config.Config, fileFacade facade.FileFacade, useCases *usecase.UseCases) {
	filesHandler := handler.NewFilesHandler(fileFacade, useCases.UpdateFileUseCase)

	uploadHanler := handler.NewUploadHandler(config, useCases.UploadUseCase, useCases.CreateFileUseCase)

	downloadHandler := handler.NewDownloadHandler(useCases.DownloadFileUseCase, fileFacade)

	router := NewFilesRouter(config, filesHandler, uploadHanler, downloadHandler).MountRoutes()
	http.Handle("/", router)
	slog.Info("File Manager REST API runing", "port", config.Server.Port)

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", config.Server.Port),
		ReadHeaderTimeout: time.Duration(config.Server.ReadHeaderTimeout) * time.Second,
		Handler:           router,
	}

	if err := srv.ListenAndServe(); err != nil {
		slog.Error("Could not start File Manager REST server", "error", err)
		os.Exit(1)
	}
}
