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

	file, error_ = os.Open(d.config.Server.Storage.Path + "/" + fileId)

	if error_ != nil {
		slog.Error("Could not open file in fs", "traceId", traceId, "error", error_)
		return
	}

	return
}
