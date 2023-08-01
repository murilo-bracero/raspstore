package usecase

import (
	"context"
	"log"
	"os"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/murilo-bracero/raspstore/file-service/internal"
)

type DownloadFileUseCase interface {
	Execute(ctx context.Context, fileId string) (file *os.File, error_ error)
}

type downloadFileUseCase struct{}

func NewDownloadFileUseCase() DownloadFileUseCase {
	return &downloadFileUseCase{}
}

func (d *downloadFileUseCase) Execute(ctx context.Context, fileId string) (file *os.File, error_ error) {
	traceId := ctx.Value(middleware.RequestIDKey).(string)

	if error_ != nil {
		log.Printf("[ERROR] - [%s]: Could not retrieve file from file-service due to error: %s", traceId, error_.Error())
		return
	}

	file, error_ = os.Open(internal.StoragePath() + "/" + fileId)

	if error_ != nil {
		log.Printf("[ERROR] - [%s]: Could not open file in fs due to error: %s", traceId, error_.Error())
		return
	}

	return
}
