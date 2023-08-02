package api

import (
	"fmt"
	"net/http"

	"github.com/murilo-bracero/raspstore/commons/pkg/logger"
	"github.com/murilo-bracero/raspstore/file-service/internal"
	"github.com/murilo-bracero/raspstore/file-service/internal/api/handler"
	"github.com/murilo-bracero/raspstore/file-service/internal/usecase"
)

func StartApiServer(luc usecase.ListFilesUseCase,
	uuc usecase.UpdateFileUseCase,
	duc usecase.DeleteFileUseCase,
	upc usecase.UploadFileUseCase,
	downloadUc usecase.DownloadFileUseCase,
	createUc usecase.CreateFileUseCase,
	getFileUc usecase.GetFileUseCase) {
	filesHandler := handler.NewFilesHandler(luc, uuc, duc)

	uploadHanler := handler.NewUploadHandler(upc, createUc)

	downloadHandler := handler.NewDownloadHandler(downloadUc, getFileUc)

	router := NewFilesRouter(filesHandler, uploadHanler, downloadHandler).MountRoutes()
	http.Handle("/", router)
	logger.Info("File Manager REST API runing on port %d", internal.RestPort())
	http.ListenAndServe(fmt.Sprintf(":%d", internal.RestPort()), router)
}
