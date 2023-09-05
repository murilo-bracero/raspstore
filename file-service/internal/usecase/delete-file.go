package usecase

import (
	"context"

	cm "github.com/go-chi/chi/v5/middleware"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/murilo-bracero/raspstore/commons/pkg/logger"
	m "github.com/murilo-bracero/raspstore/file-service/internal/api/middleware"
	"github.com/murilo-bracero/raspstore/file-service/internal/repository"
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

	if err := u.repo.Delete(user.Subject(), fileId); err != nil {
		logger.Error("[%s]: Could not delete file in database: fileId=%s, %s", traceId, fileId, err.Error())
		return
	}

	logger.Info("[%s]: File removed successfully: fileId=%s", traceId, fileId)
	return
}
