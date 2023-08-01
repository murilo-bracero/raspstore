package main

import (
	"context"
	"os"

	"github.com/joho/godotenv"
	"github.com/murilo-bracero/raspstore/commons/pkg/logger"
	"github.com/murilo-bracero/raspstore/file-service/internal/api"
	db "github.com/murilo-bracero/raspstore/file-service/internal/database"
	"github.com/murilo-bracero/raspstore/file-service/internal/repository"
	"github.com/murilo-bracero/raspstore/file-service/internal/usecase"
)

func main() {
	ctx := context.Background()

	if err := godotenv.Load(); err != nil {
		logger.Warn("Could not load .env file. Using system variables instead")
	}

	conn, err := db.NewMongoConnection(context.Background())

	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	defer conn.Close(context.Background())

	fileRepo := repository.NewFilesRepository(ctx, conn)

	useCases := usecase.InitUseCases(fileRepo)

	logger.Info("Bootstraping servers")
	api.StartApiServer(useCases.ListFilesUseCase,
		useCases.UpdateFileUseCase,
		useCases.DeleteFileUseCase,
		useCases.UploadUseCase,
		useCases.DownloadFileUseCase,
		useCases.CreateFileUseCase,
		useCases.GetFileUseCase)
}
