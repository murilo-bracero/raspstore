package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"raspstore.github.io/fs-service/api"
	"raspstore.github.io/fs-service/api/controller"
	"raspstore.github.io/fs-service/internal"
	"raspstore.github.io/fs-service/usecase"
)

func main() {
	ctx := context.Background()

	if err := godotenv.Load(); err != nil {
		log.Println("Could not load .env file. Using system variables instead")
	}

	if _, err := os.Stat(internal.StoragePath()); os.IsNotExist(err) {
		log.Println("[INFO] Base path does not exists, starting creation")
		err = os.MkdirAll(internal.StoragePath(), 0755)

		if err != nil {
			log.Println("[ERROR] Could not create base path: ", err.Error())
		}
	}

	fileUseCase := usecase.NewFileInfoUseCase(ctx)

	ctr := controller.NewFileServeController(fileUseCase)

	router := api.NewRoutes(ctr).MountRoutes()

	http.Handle("/", router)
	log.Printf("File Manager API runing on port %d", internal.RestPort())
	http.ListenAndServe(fmt.Sprintf(":%d", internal.RestPort()), router)
}
