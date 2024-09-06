package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/murilo-bracero/raspstore/file-service/internal/application/facade"
	"github.com/murilo-bracero/raspstore/file-service/internal/application/usecase"
	"github.com/murilo-bracero/raspstore/file-service/internal/infra/config"
	"github.com/murilo-bracero/raspstore/file-service/internal/infra/handler"
	"github.com/murilo-bracero/raspstore/file-service/internal/infra/validator"
)

func StartApiServer(config *config.Config, fileFacade facade.FileFacade, useCases *usecase.UseCases) {
	appHandler := handler.New(useCases.DownloadFileUseCase,
		useCases.UploadUseCase,
		useCases.CreateFileUseCase,
		useCases.UpdateFileUseCase,
		useCases.LoginPAMUseCase,
		fileFacade,
		config)

	jwtValidator, err := validator.NewJWTValidator(context.Background(), config)

	if err != nil {
		slog.Error("Could not setup JWT validator", "err", err)
		os.Exit(1)
	}

	router := NewFilesRouter(config, appHandler, jwtValidator).MountRoutes()
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
