package usecase

import (
	"context"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/murilo-bracero/raspstore/commons/pkg/logger"
	"github.com/murilo-bracero/raspstore/idp/internal/model"
	"github.com/murilo-bracero/raspstore/idp/internal/repository"
)

type GetUserUseCase interface {
	Execute(ctx context.Context, userId string) (user *model.User, error_ error)
}

type getProfileUseCase struct {
	userRepository repository.UsersRepository
}

func NewGetUserUseCase(userRepository repository.UsersRepository) GetUserUseCase {
	return &getProfileUseCase{userRepository: userRepository}
}

func (u *getProfileUseCase) Execute(ctx context.Context, userId string) (user *model.User, error_ error) {
	traceId := ctx.Value(middleware.RequestIDKey).(string)
	user, error_ = u.userRepository.FindByUserId(userId)

	if error_ != nil {
		logger.Error("[%s] Could not search for user: userId=%s : %s", traceId, userId, error_.Error())
	}

	return
}