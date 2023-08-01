package main

import (
	"context"
	"os"
	"sync"

	"github.com/joho/godotenv"
	"github.com/murilo-bracero/raspstore/commons/pkg/logger"
	"github.com/murilo-bracero/raspstore/file-info-service/internal/api"
	db "github.com/murilo-bracero/raspstore/file-info-service/internal/database"
	"github.com/murilo-bracero/raspstore/file-info-service/internal/grpc"
	"github.com/murilo-bracero/raspstore/file-info-service/internal/repository"
	"github.com/murilo-bracero/raspstore/file-info-service/internal/usecase"
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

	var wg sync.WaitGroup

	wg.Add(2)
	logger.Info("Bootstraping servers")
	go grpc.StartGrpcServer(useCases.GetFileUseCase, useCases.CreateFileUseCase)
	go api.StartApiServer(useCases.ListFilesUseCase, useCases.UpdateFileUseCase, useCases.DeleteFileUseCase, useCases.UploadUseCase, useCases.CreateFileUseCase)
	wg.Wait()
}
