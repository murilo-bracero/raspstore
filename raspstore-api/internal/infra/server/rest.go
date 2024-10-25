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

type ApiServerParams struct {
	Config            *config.Config
	FileFacade        facade.FileFacade
	FileSystemFacade  facade.FileSystemFacade
	UpdateFileUseCase usecase.UpdateFileUseCase
}

func StartApiServer(params *ApiServerParams) {
	appHandler := handler.New(
		params.UpdateFileUseCase,
		params.FileFacade,
		params.FileSystemFacade,
		params.Config)

	jwtValidator, err := validator.NewJWTValidator(context.Background(), params.Config)

	if err != nil {
		slog.Error("Could not setup JWT validator", "err", err)
		os.Exit(1)
	}

	router := NewFilesRouter(params.Config, appHandler, jwtValidator).MountRoutes()
	http.Handle("/", router)
	slog.Info("File Manager REST API runing", "port", params.Config.Server.Port)

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", params.Config.Server.Port),
		ReadHeaderTimeout: time.Duration(params.Config.Server.ReadHeaderTimeout) * time.Second,
		Handler:           router,
	}

	if err := srv.ListenAndServe(); err != nil {
		slog.Error("Could not start File Manager REST server", "error", err)
		os.Exit(1)
	}
}
