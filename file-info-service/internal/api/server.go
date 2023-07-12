package api

import (
	"fmt"
	"log"
	"net/http"

	"raspstore.github.io/file-manager/internal"
	"raspstore.github.io/file-manager/internal/api/handler"
	"raspstore.github.io/file-manager/internal/grpc"
	"raspstore.github.io/file-manager/internal/repository"
)

func StartApiServer(fileRepository repository.FilesRepository) {
	filesHandler := handler.NewFilesHandler(fileRepository)

	authService := grpc.NewAuthService()

	router := NewFilesRouter(filesHandler, authService).MountRoutes()
	http.Handle("/", router)
	log.Printf("File Manager REST API runing on port %d", internal.RestPort())
	http.ListenAndServe(fmt.Sprintf(":%d", internal.RestPort()), router)
}
