package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/murilo-bracero/raspstore/file-info-service/internal"
	"github.com/murilo-bracero/raspstore/file-info-service/internal/api/handler"
	"github.com/murilo-bracero/raspstore/file-info-service/internal/grpc/client"
	"github.com/murilo-bracero/raspstore/file-info-service/internal/usecase"
)

func StartApiServer(luc usecase.ListFilesUseCase, uuc usecase.UpdateFileUseCase, duc usecase.DeleteFileUseCase) {
	filesHandler := handler.NewFilesHandler(luc, uuc, duc)

	authService := client.NewAuthService()

	router := NewFilesRouter(filesHandler, authService).MountRoutes()
	http.Handle("/", router)
	log.Printf("File Manager REST API runing on port %d", internal.RestPort())
	http.ListenAndServe(fmt.Sprintf(":%d", internal.RestPort()), router)
}
