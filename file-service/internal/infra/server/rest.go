package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

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

	if err := http.ListenAndServe(fmt.Sprintf(":%d", config.Server.Port), router); err != nil {
		slog.Error("Could not start File Manager REST server", "error", err)
		os.Exit(1)
	}
}
