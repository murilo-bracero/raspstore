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
	Execute(ctx context.Context, file *model.File) (fileMetadata *model.File, error_ error)
}

type updateFileUseCase struct {
	repo repository.FilesRepository
}

func NewUpdateFileUseCase(repo repository.FilesRepository) UpdateFileUseCase {
	return &updateFileUseCase{repo: repo}
}

func (c *updateFileUseCase) Execute(ctx context.Context, file *model.File) (fileMetadata *model.File, error_ error) {
	user := ctx.Value(m.UserClaimsCtxKey).(jwt.Token)
	traceId := ctx.Value(chiMiddleware.RequestIDKey).(string)

	found, error_ := c.repo.FindById(user.Subject(), file.FileId)

	if error_ != nil {
		slog.Error("Could not find file", "traceId", traceId, "fileId", file.FileId, "error", error_)
		return
	}

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
		slog.Error("Could not update file", "traceId", traceId, "fileId", file.FileId, "error", error_)
		return
	}

	slog.Info("File updated successfully", "traceId", traceId, "fileId", file.FileId)

	fileMetadata, error_ = c.repo.FindByIdLookup(user.Subject(), file.FileId)

	if error_ != nil {
		slog.Error("Could not search lookup file", "traceId", traceId, "fileId", file.FileId, "error", error_)
		return
	}

	return
}
