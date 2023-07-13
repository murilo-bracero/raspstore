package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"raspstore.github.io/fs-service/internal"
	"raspstore.github.io/fs-service/internal/api"
	"raspstore.github.io/fs-service/internal/grpc/client"
	"raspstore.github.io/fs-service/internal/usecase"
)

func main() {
	ctx := context.Background()

	if err := godotenv.Load(); err != nil {
		log.Println("Could not load .env file. Using system variables instead")
	}

	checkStoragePath()

	authService := client.NewAuthService(ctx)
	fileInfoService := client.NewFileInfoService(ctx)

	uploadUseCase := usecase.NewUploadFileUseCase(fileInfoService)
	downloadUseCase := usecase.NewDownloadFileUseCase(fileInfoService)

	api.StartRestServer(uploadUseCase, downloadUseCase, authService)
}

func checkStoragePath() {
	if _, err := os.Stat(internal.StoragePath()); os.IsNotExist(err) {
		log.Println("[INFO] Base path does not exists, starting creation")
		err = os.MkdirAll(internal.StoragePath(), 0755)

		if err != nil {
			log.Println("[ERROR] Could not create base path: ", err.Error())
		}
	}
}
