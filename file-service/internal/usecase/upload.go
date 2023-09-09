package usecase

import (
	"context"
	"io"
	"log/slog"
	"os"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/murilo-bracero/raspstore/file-service/internal/infra"
	"github.com/murilo-bracero/raspstore/file-service/internal/model"
)

type UploadFileUseCase interface {
	Execute(ctx context.Context, file *model.File, src io.Reader) (error_ error)
}

type uploadFileUseCase struct {
	config *infra.Config
}

func NewUploadFileUseCase(config *infra.Config) UploadFileUseCase {
	return &uploadFileUseCase{config: config}
}

func (u *uploadFileUseCase) Execute(ctx context.Context, file *model.File, src io.Reader) (error_ error) {
	traceId := ctx.Value(middleware.RequestIDKey).(string)

	filerep, error_ := os.Create(u.config.Server.Storage.Path + "/" + file.FileId)

	if error_ != nil {
		slog.Error("[%s]: Could not create file in fs due to error: %s", traceId, error_.Error())
		return
	}

	defer filerep.Close()

	_, error_ = io.Copy(filerep, src)

	if error_ != nil {
		slog.Error("Could not read file buffer due to error", "traceId", traceId, "error", error_)
		return
	}

	return
}
