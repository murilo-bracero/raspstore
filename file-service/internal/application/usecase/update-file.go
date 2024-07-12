package usecase

import (
	"context"
	"log/slog"

	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/murilo-bracero/raspstore/file-service/internal/application/repository"
	"github.com/murilo-bracero/raspstore/file-service/internal/domain/entity"
	m "github.com/murilo-bracero/raspstore/file-service/internal/infra/middleware"
)

type UpdateFileUseCase interface {
	Execute(ctx context.Context, file *entity.File) (fileMetadata *entity.File, err error)
}

type updateFileUseCase struct {
	repo repository.FilesRepository
}

func NewUpdateFileUseCase(repo repository.FilesRepository) UpdateFileUseCase {
	return &updateFileUseCase{repo: repo}
}

func (c *updateFileUseCase) Execute(ctx context.Context, file *entity.File) (fileMetadata *entity.File, err error) {
	user := ctx.Value(m.UserClaimsCtxKey).(jwt.Token)
	traceId := ctx.Value(chiMiddleware.RequestIDKey).(string)

	fileMetadata, err = c.repo.FindById(user.Subject(), file.FileId)

	if err != nil {
		slog.Error("Could not find file", "traceId", traceId, "fileId", file.FileId, "error", err)
		return
	}

	fileMetadata.Secret = file.Secret
	fileMetadata.Filename = file.Filename

	if fileMetadata.Secret {
		fileMetadata.Viewers = []string{}
		fileMetadata.Editors = []string{}
	} else {
		fileMetadata.Viewers = file.Viewers
		fileMetadata.Editors = file.Editors
	}

	if err = c.repo.Update(user.Subject(), fileMetadata); err != nil {
		slog.Error("Could not update file", "traceId", traceId, "fileId", file.FileId, "error", err)
		return nil, err
	}

	slog.Info("File updated successfully", "traceId", traceId, "fileId", file.FileId)

	return
}
