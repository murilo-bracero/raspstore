package usecase

import (
	"context"
	"io"
	"log"
	"os"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/murilo-bracero/raspstore/file-info-service/internal"
	"github.com/murilo-bracero/raspstore/file-info-service/internal/model"
)

type UploadFileUseCase interface {
	Execute(ctx context.Context, file *model.File, src io.Reader) (error_ error)
}

type uploadFileUseCase struct {
}

func NewUploadFileUseCase() UploadFileUseCase {
	return &uploadFileUseCase{}
}

func (u *uploadFileUseCase) Execute(ctx context.Context, file *model.File, src io.Reader) (error_ error) {
	traceId := ctx.Value(middleware.RequestIDKey).(string)

	filerep, error_ := os.Create(internal.StoragePath() + "/" + file.FileId)

	if error_ != nil {
		log.Printf("[ERROR] - [%s]: Could not create file in fs due to error: %s", traceId, error_.Error())
		return
	}

	defer filerep.Close()

	_, error_ = io.Copy(filerep, src)

	if error_ != nil {
		log.Printf("[ERROR] - [%s]: Could not read file buffer due to error: %s", traceId, error_.Error())
		return
	}

	return
}
