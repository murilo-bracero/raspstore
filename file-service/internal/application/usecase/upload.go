package usecase

import (
	"context"
	"io"
	"log/slog"
	"os"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/murilo-bracero/raspstore/file-service/internal/domain/entity"
	"github.com/murilo-bracero/raspstore/file-service/internal/infra/config"
)

type UploadFileUseCase interface {
	Execute(ctx context.Context, file *entity.File, src io.Reader) (err error)
}

type uploadFileUseCase struct {
	config *config.Config
}

func NewUploadFileUseCase(config *config.Config) *uploadFileUseCase {
	return &uploadFileUseCase{config: config}
}

func (u *uploadFileUseCase) Execute(ctx context.Context, file *entity.File, src io.Reader) (err error) {
	traceId := ctx.Value(middleware.RequestIDKey).(string)

	filerep, err := os.Create(u.config.Server.Storage.Path + "/" + file.FileId)

	if err != nil {
		slog.Error("Could not create file in fs", "traceId", traceId, "error", err)
		return
	}

	defer filerep.Close()

	_, err = io.Copy(filerep, src)

	if err != nil {
		slog.Error("Could not read file buffer", "traceId", traceId, "error", err)
		return
	}

	return
}
