package usecase

import (
	"context"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/murilo-bracero/raspstore/commons/pkg/logger"
	"github.com/murilo-bracero/raspstore/idp/internal/repository"
)

type DeleteUserUseCase interface {
	Execute(ctx context.Context, userId string) error
}

type deleteUserUseCase struct {
	userRepository repository.UsersRepository
}

func NewDeleteUserUseCase(userRepository repository.UsersRepository) DeleteUserUseCase {
	return &deleteUserUseCase{userRepository: userRepository}
}

func (u *deleteUserUseCase) Execute(ctx context.Context, userId string) error {
	err := u.userRepository.Delete(userId)

	if err != nil {
		traceId := ctx.Value(middleware.RequestIDKey).(string)
		logger.Error("[%s] Could not remove user: userId=%s : %s", traceId, userId, err.Error())
	}

	return err
}
