package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/murilo-bracero/raspstore/commons/pkg/service"
	"raspstore.github.io/fs-service/internal"
	"raspstore.github.io/fs-service/internal/api/handler"
	"raspstore.github.io/fs-service/internal/usecase"
)

func StartRestServer(uuc usecase.UploadFileUseCase, duc usecase.DownloadFileUseCase, as service.AuthService) {
	hndlr := handler.NewFileServeHandler(uuc, duc)

	router := NewRoutes(hndlr, as).MountRoutes()

	http.Handle("/", router)
	log.Printf("File Manager API runing on port %d", internal.RestPort())
	http.ListenAndServe(fmt.Sprintf(":%d", internal.RestPort()), router)
}
