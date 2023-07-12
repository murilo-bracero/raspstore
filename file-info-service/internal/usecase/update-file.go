package usecase

import (
	"context"
	"log"

	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	rMiddleware "github.com/murilo-bracero/raspstore/commons/pkg/middleware"
	"github.com/murilo-bracero/raspstore/file-info-service/internal/model"
	"github.com/murilo-bracero/raspstore/file-info-service/internal/repository"
)

type UpdateFileUseCase interface {
	Execute(ctx context.Context, file *model.File) (fileMetadata *model.FileMetadataLookup, error_ error)
}

type updateFileUseCase struct {
	repo repository.FilesRepository
}

func NewUpdateFileUseCase(repo repository.FilesRepository) UpdateFileUseCase {
	return &updateFileUseCase{repo: repo}
}

func (c *updateFileUseCase) Execute(ctx context.Context, file *model.File) (fileMetadata *model.FileMetadataLookup, error_ error) {
	userId := ctx.Value(rMiddleware.UserIdKey).(string)
	traceId := ctx.Value(chiMiddleware.RequestIDKey).(string)

	found, error_ := c.repo.FindById(userId, file.FileId)

	if error_ != nil {
		log.Printf("[ERROR] - [%s]: Could not search file with id %s in database: %s", traceId, file.FileId, error_.Error())
		return
	}

	log.Printf("[INFO] - [%s]: File with id=%s found", traceId, file.FileId)

	if file.Filename != "" {
		found.Filename = file.Filename
	}

	if file.Path != "" {
		found.Path = file.Path
	}

	if file.Viewers != nil || len(file.Viewers) != 0 {
		found.Viewers = file.Viewers
	}

	if file.Editors != nil || len(file.Editors) != 0 {
		found.Editors = file.Editors
	}

	if error_ = c.repo.Update(userId, found); error_ != nil {
		log.Printf("[ERROR] - [%s]: Could not update file with id %s in database: %s", traceId, file.FileId, error_.Error())
		return
	}

	log.Printf("[INFO] - [%s]: File with id=%s updated successfully", traceId, file.FileId)

	fileMetadata, error_ = c.repo.FindByIdLookup(userId, file.FileId)

	if error_ != nil {
		log.Printf("[ERROR] - [%s]: Could not search lookup file with id %s in database: %s", traceId, file.FileId, error_.Error())
		return
	}

	return
}
