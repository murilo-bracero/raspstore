package main

import (
	"context"
	"log"
	"sync"

	"github.com/joho/godotenv"
	"github.com/murilo-bracero/raspstore/file-info-service/internal/api"
	db "github.com/murilo-bracero/raspstore/file-info-service/internal/database"
	"github.com/murilo-bracero/raspstore/file-info-service/internal/grpc"
	"github.com/murilo-bracero/raspstore/file-info-service/internal/grpc/client"
	"github.com/murilo-bracero/raspstore/file-info-service/internal/repository"
	"github.com/murilo-bracero/raspstore/file-info-service/internal/usecase"
)

func main() {
	ctx := context.Background()

	if err := godotenv.Load(); err != nil {
		log.Println("Could not load .env file. Using system variables instead")
	}

	conn, err := db.NewMongoConnection(context.Background())

	if err != nil {
		log.Panicln(err)
	}

	defer conn.Close(context.Background())

	fileRepo := repository.NewFilesRepository(ctx, conn)

	userServiceClient := client.NewUserConfigGrpcService()

	useCases := usecase.InitUseCases(fileRepo, userServiceClient)

	var wg sync.WaitGroup

	wg.Add(2)
	log.Println("[INFO] Bootstraping servers")
	go grpc.StartGrpcServer(useCases.GetFileUseCase, useCases.CreateFileUseCase)
	go api.StartApiServer(useCases.ListFilesUseCase, useCases.UpdateFileUseCase, useCases.DeleteFileUseCase)
	wg.Wait()
}
