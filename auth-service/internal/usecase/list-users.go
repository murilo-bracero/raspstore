package usecase

import (
	"context"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/murilo-bracero/raspstore/auth-service/internal/model"
	"github.com/murilo-bracero/raspstore/auth-service/internal/repository"
	"github.com/murilo-bracero/raspstore/commons/pkg/logger"
)

const defaultMaxSize = 50

type ListUsersUseCase interface {
	Execute(ctx context.Context, page int, size int, username string, enabled *bool) (userPage *model.UserPage, error_ error)
}

type listUsersUseCase struct {
	userRepository repository.UsersRepository
}

func NewListUsersUseCase(userRepository repository.UsersRepository) ListUsersUseCase {
	return &listUsersUseCase{userRepository: userRepository}
}

func (u *listUsersUseCase) Execute(ctx context.Context, page int, size int, username string, enabled *bool) (userPage *model.UserPage, error_ error) {
	if size == 0 || size > defaultMaxSize {
		size = defaultMaxSize
	}

	userPage, error_ = u.userRepository.FindAll(page, size, username, enabled)

	if error_ != nil {
		traceId := ctx.Value(middleware.RequestIDKey).(string)
		logger.Error("[%s] Could not list users due to error: %s", traceId, error_.Error())
	}

	return
}
