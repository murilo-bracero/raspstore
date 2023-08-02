package usecase

import (
	"context"

	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/murilo-bracero/raspstore/commons/pkg/logger"
	rmd "github.com/murilo-bracero/raspstore/commons/pkg/security/middleware"
	"github.com/murilo-bracero/raspstore/file-service/internal/model"
	"github.com/murilo-bracero/raspstore/file-service/internal/repository"
)

type ListFilesUseCase interface {
	Execute(ctx context.Context, page int, size int, filename string, secret bool) (filesPage *model.FilePage, error_ error)
}

const maxListSize = 50

type listFilesUseCase struct {
	repo repository.FilesRepository
}

func NewListFilesUseCase(repo repository.FilesRepository) ListFilesUseCase {
	return &listFilesUseCase{repo: repo}
}

func (u *listFilesUseCase) Execute(ctx context.Context, page int, size int, filename string, secret bool) (filesPage *model.FilePage, error_ error) {
	if size == 0 || size > maxListSize {
		size = maxListSize
	}

	userId := ctx.Value(rmd.UserClaimsCtxKey).(string)

	filesPage, error_ = u.repo.FindAll(userId, page, size, filename, secret)

	if error_ != nil {
		traceId := ctx.Value(chiMiddleware.RequestIDKey).(string)
		logger.Error("[%s]: Could list files due to error: %s", traceId, error_.Error())
		return
	}

	return
}
