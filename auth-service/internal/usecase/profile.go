package usecase

import (
	"context"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/murilo-bracero/raspstore/auth-service/internal/model"
	"github.com/murilo-bracero/raspstore/auth-service/internal/repository"
	"github.com/murilo-bracero/raspstore/commons/pkg/logger"
)

type GetProfileUseCase interface {
	Execute(ctx context.Context, userId string) (user *model.User, error_ error)
}

type getProfileUseCase struct {
	userRepository repository.UsersRepository
}

func NewGetProfileUseCase(userRepository repository.UsersRepository) GetProfileUseCase {
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
