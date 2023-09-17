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

type ListFilesUseCase interface {
	Execute(ctx context.Context, page int, size int, filename string, secret bool) (filesPage *entity.FilePage, error_ error)
}

const maxListSize = 50

type listFilesUseCase struct {
	repo repository.FilesRepository
}

func NewListFilesUseCase(repo repository.FilesRepository) ListFilesUseCase {
	return &listFilesUseCase{repo: repo}
}

func (u *listFilesUseCase) Execute(ctx context.Context, page int, size int, filename string, secret bool) (filesPage *entity.FilePage, error_ error) {
	if size == 0 || size > maxListSize {
		size = maxListSize
	}

	user := ctx.Value(m.UserClaimsCtxKey).(jwt.Token)

	filesPage, error_ = u.repo.FindAll(user.Subject(), page, size, filename, secret)

	if error_ != nil {
		traceId := ctx.Value(chiMiddleware.RequestIDKey).(string)
		slog.Error("Could list files due to error", "traceId", traceId, "error", error_)
		return
	}

	return
}
