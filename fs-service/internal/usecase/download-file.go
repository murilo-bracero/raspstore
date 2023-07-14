package usecase

import (
	"context"
	"log"
	"os"

	"github.com/go-chi/chi/v5/middleware"
	rMiddleware "github.com/murilo-bracero/raspstore/commons/pkg/middleware"
	"raspstore.github.io/fs-service/internal"
	"raspstore.github.io/fs-service/internal/grpc/client"
	"raspstore.github.io/fs-service/internal/model"
)

type DownloadFileUseCase interface {
	Execute(ctx context.Context, fileId string) (downloadRep *model.FileDownloadRepresentation, error_ error)
}

type downloadFileUseCase struct {
	service client.FileInfoService
}

func NewDownloadFileUseCase(service client.FileInfoService) DownloadFileUseCase {
	return &downloadFileUseCase{service: service}
}

func (d *downloadFileUseCase) Execute(ctx context.Context, fileId string) (downloadRep *model.FileDownloadRepresentation, error_ error) {
	traceId := ctx.Value(middleware.RequestIDKey).(string)
	userId := ctx.Value(rMiddleware.UserIdKey).(string)

	fileInfo, error_ := d.service.GetFileMetadataById(fileId, userId)

	if error_ != nil {
		log.Printf("[ERROR] - [%s]: Could not retrieve file from file-info-service due to error: %s", traceId, error_.Error())
		return
	}

	file, error_ := os.Open(internal.StoragePath() + "/" + fileInfo.FileId)

	if error_ != nil {
		log.Printf("[ERROR] - [%s]: Could not open file in fs due to error: %s", traceId, error_.Error())
		return
	}

	fileStat, error_ := file.Stat()

	if error_ != nil {
		log.Printf("[ERROR] - [%s]: Could not get status from fs due to error: %s", traceId, error_.Error())
		return
	}

	return &model.FileDownloadRepresentation{
		Filename: fileInfo.Filename,
		File:     file,
		FileSize: fileStat.Size(),
	}, nil
}
