package usecase

import (
	"context"
	"log/slog"
	"os"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/murilo-bracero/raspstore/file-service/internal/infra/config"
)

type DownloadFileUseCase interface {
	Execute(ctx context.Context, fileId string) (file *os.File, error_ error)
}

type downloadFileUseCase struct {
	config *config.Config
}

func NewDownloadFileUseCase(config *config.Config) DownloadFileUseCase {
	return &downloadFileUseCase{config: config}
}

func (d *downloadFileUseCase) Execute(ctx context.Context, fileId string) (file *os.File, error_ error) {
	traceId := ctx.Value(middleware.RequestIDKey).(string)

	if error_ != nil {
		slog.Error("[%s]: Could not retrieve file from file-service due to error: %s", traceId, error_.Error())
		return
	}

	file, error_ = os.Open(d.config.Server.Storage.Path + "/" + fileId)

	if error_ != nil {
		slog.Error("[%s]: Could not open file in fs due to error: %s", traceId, error_.Error())
		return
	}

	return
}
