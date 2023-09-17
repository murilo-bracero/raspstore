package usecase

import (
	"context"
	"log/slog"

	cm "github.com/go-chi/chi/v5/middleware"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/murilo-bracero/raspstore/file-service/internal/application/repository"
	m "github.com/murilo-bracero/raspstore/file-service/internal/infra/middleware"
)

type DeleteFileUseCase interface {
	Execute(ctx context.Context, fileId string) (error_ error)
}

type deleteFileUseCase struct {
	repo repository.FilesRepository
}

func NewDeleteFileUseCase(repo repository.FilesRepository) DeleteFileUseCase {
	return &deleteFileUseCase{repo: repo}
}

func (u *deleteFileUseCase) Execute(ctx context.Context, fileId string) (error_ error) {
	traceId := ctx.Value(cm.RequestIDKey).(string)
	user := ctx.Value(m.UserClaimsCtxKey).(jwt.Token)

	if error_ = u.repo.Delete(user.Subject(), fileId); error_ != nil {
		slog.Error("Could not delete file in database:", "traceId", traceId, "fileId", fileId, "error", error_)
		return
	}

	slog.Info("File removed successfully:", "traceId", traceId, "fileId", fileId)
	return
}
