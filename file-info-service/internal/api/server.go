package api

import (
	"fmt"
	"net/http"

	"github.com/murilo-bracero/raspstore/commons/pkg/logger"
	"github.com/murilo-bracero/raspstore/file-info-service/internal"
	"github.com/murilo-bracero/raspstore/file-info-service/internal/api/handler"
	"github.com/murilo-bracero/raspstore/file-info-service/internal/grpc/client"
	"github.com/murilo-bracero/raspstore/file-info-service/internal/usecase"
)

func StartApiServer(luc usecase.ListFilesUseCase,
	uuc usecase.UpdateFileUseCase,
	duc usecase.DeleteFileUseCase,
	upc usecase.UploadFileUseCase,
	createUc usecase.CreateFileUseCase) {
	filesHandler := handler.NewFilesHandler(luc, uuc, duc)

	uploadHanler := handler.NewUploadHandler(upc, createUc)

	authService := client.NewAuthService()

	router := NewFilesRouter(filesHandler, uploadHanler, authService).MountRoutes()
	http.Handle("/", router)
	logger.Info("File Manager REST API runing on port %d", internal.RestPort())
	http.ListenAndServe(fmt.Sprintf(":%d", internal.RestPort()), router)
}
