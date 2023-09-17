package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"github.com/murilo-bracero/raspstore/file-service/internal/application/usecase"
	"github.com/murilo-bracero/raspstore/file-service/internal/infra/config"
	"github.com/murilo-bracero/raspstore/file-service/internal/infra/repository"
	"github.com/murilo-bracero/raspstore/file-service/internal/infra/server"
)

func main() {

	ctx := context.Background()

	if err := godotenv.Load(); err != nil {
		slog.Warn("Could not load .env file. Using system variables instead")
	}

	config := config.NewConfig("config/application.yaml")

	conn, err := repository.NewDatabaseConnection(ctx, config)

	if err != nil {
		slog.Error("could not connect to database", "error", err)
		os.Exit(1)
	}

	defer conn.Close(ctx)

	fileRepo := repository.NewFilesRepository(ctx, conn)

	useCases := usecase.InitUseCases(config, fileRepo)

	slog.Info("Bootstraping servers")
	server.StartApiServer(config,
		useCases.ListFilesUseCase,
		useCases.UpdateFileUseCase,
		useCases.DeleteFileUseCase,
		useCases.UploadUseCase,
		useCases.DownloadFileUseCase,
		useCases.CreateFileUseCase,
		useCases.GetFileUseCase)
}
