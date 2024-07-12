package usecase

import (
	"context"
	"log/slog"
	"os"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/murilo-bracero/raspstore/file-service/internal/infra/config"
)

type DownloadFileUseCase interface {
	Execute(ctx context.Context, fileId string) (file *os.File, err error)
}

type downloadFileUseCase struct {
	config *config.Config
}

func NewDownloadFileUseCase(config *config.Config) DownloadFileUseCase {
	return &downloadFileUseCase{config: config}
}

func (d *downloadFileUseCase) Execute(ctx context.Context, fileId string) (file *os.File, err error) {
	traceId := ctx.Value(middleware.RequestIDKey).(string)

	file, err = os.Open(d.config.Server.Storage.Path + "/" + fileId)

	if err != nil {
		slog.Error("Could not open file in fs", "traceId", traceId, "error", err)
		return
	}

	return
}
