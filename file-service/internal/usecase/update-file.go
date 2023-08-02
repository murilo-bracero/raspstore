package usecase

import (
	"context"

	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/murilo-bracero/raspstore/commons/pkg/logger"
	rmd "github.com/murilo-bracero/raspstore/commons/pkg/security/middleware"
	"github.com/murilo-bracero/raspstore/file-service/internal/model"
	"github.com/murilo-bracero/raspstore/file-service/internal/repository"
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
	userId := ctx.Value(rmd.UserClaimsCtxKey).(string)
	traceId := ctx.Value(chiMiddleware.RequestIDKey).(string)

	found, error_ := c.repo.FindById(userId, file.FileId)

	if error_ != nil {
		logger.Error("[%s]: Could not search file with id %s in database: %s", traceId, file.FileId, error_.Error())
		return
	}

	logger.Info("[%s]: File with id=%s found", traceId, file.FileId)

	found.Secret = file.Secret
	found.Filename = file.Filename

	if found.Secret {
		found.Viewers = []string{}
		found.Editors = []string{}
	} else {
		found.Viewers = file.Viewers
		found.Editors = file.Editors
	}

	if error_ = c.repo.Update(userId, found); error_ != nil {
		logger.Error("[%s]: Could not update file with id %s in database: %s", traceId, file.FileId, error_.Error())
		return
	}

	logger.Info("[%s]: File with id=%s updated successfully", traceId, file.FileId)

	fileMetadata, error_ = c.repo.FindByIdLookup(userId, file.FileId)

	if error_ != nil {
		logger.Error("[%s]: Could not search lookup file with id %s in database: %s", traceId, file.FileId, error_.Error())
		return
	}

	return
}
