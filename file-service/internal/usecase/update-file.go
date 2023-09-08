package usecase

import (
	"context"
	"log/slog"

	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/lestrrat-go/jwx/jwt"
	m "github.com/murilo-bracero/raspstore/file-service/internal/api/middleware"
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
	user := ctx.Value(m.UserClaimsCtxKey).(jwt.Token)
	traceId := ctx.Value(chiMiddleware.RequestIDKey).(string)

	found, error_ := c.repo.FindById(user.Subject(), file.FileId)

	if error_ != nil {
		slog.Error("[%s]: Could not search file with id %s in database: %s", traceId, file.FileId, error_.Error())
		return
	}

	slog.Info("[%s]: File with id=%s found", traceId, file.FileId)

	found.Secret = file.Secret
	found.Filename = file.Filename

	if found.Secret {
		found.Viewers = []string{}
		found.Editors = []string{}
	} else {
		found.Viewers = file.Viewers
		found.Editors = file.Editors
	}

	if error_ = c.repo.Update(user.Subject(), found); error_ != nil {
		slog.Error("[%s]: Could not update file with id %s in database: %s", traceId, file.FileId, error_.Error())
		return
	}

	slog.Info("[%s]: File with id=%s updated successfully", traceId, file.FileId)

	fileMetadata, error_ = c.repo.FindByIdLookup(user.Subject(), file.FileId)

	if error_ != nil {
		slog.Error("[%s]: Could not search lookup file with id %s in database: %s", traceId, file.FileId, error_.Error())
		return
	}

	return
}
