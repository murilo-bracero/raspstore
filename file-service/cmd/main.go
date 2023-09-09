package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"github.com/murilo-bracero/raspstore/file-service/internal/api"
	db "github.com/murilo-bracero/raspstore/file-service/internal/database"
	"github.com/murilo-bracero/raspstore/file-service/internal/infra"
	"github.com/murilo-bracero/raspstore/file-service/internal/repository"
	"github.com/murilo-bracero/raspstore/file-service/internal/usecase"
)

func main() {

	ctx := context.Background()

	if err := godotenv.Load(); err != nil {
		slog.Warn("Could not load .env file. Using system variables instead")
	}

	config := infra.NewConfig("")

	conn, err := db.NewMongoConnection(ctx, config)

	if err != nil {
		slog.Error("could not connect to database", "error", err)
		os.Exit(1)
	}

	defer conn.Close(ctx)

	fileRepo := repository.NewFilesRepository(ctx, conn)

	useCases := usecase.InitUseCases(config, fileRepo)

	slog.Info("Bootstraping servers")
	api.StartApiServer(config,
		useCases.ListFilesUseCase,
		useCases.UpdateFileUseCase,
		useCases.DeleteFileUseCase,
		useCases.UploadUseCase,
		useCases.DownloadFileUseCase,
		useCases.CreateFileUseCase,
		useCases.GetFileUseCase)
}
