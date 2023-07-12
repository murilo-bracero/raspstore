package usecase

import (
	"context"
	"log"

	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	rMiddleware "github.com/murilo-bracero/raspstore/commons/pkg/middleware"
	"github.com/murilo-bracero/raspstore/file-info-service/internal/model"
	"github.com/murilo-bracero/raspstore/file-info-service/internal/repository"
)

type ListFilesUseCase interface {
	Execute(ctx context.Context, page int, size int) (filesPage *model.FilePage, error_ error)
}

const maxListSize = 50

type listFilesUseCase struct {
	repo repository.FilesRepository
}

func NewListFilesUseCase(repo repository.FilesRepository) ListFilesUseCase {
	return &listFilesUseCase{repo: repo}
}

func (u *listFilesUseCase) Execute(ctx context.Context, page int, size int) (filesPage *model.FilePage, error_ error) {
	if size == 0 || size > maxListSize {
		size = maxListSize
	}

	userId := ctx.Value(rMiddleware.UserIdKey).(string)

	filesPage, error_ = u.repo.FindAll(userId, page, size)

	if error_ != nil {
		traceId := ctx.Value(chiMiddleware.RequestIDKey).(string)
		log.Printf("[ERROR] - [%s]: Could list files due to error: %s", traceId, error_.Error())
		return
	}

	return
}
