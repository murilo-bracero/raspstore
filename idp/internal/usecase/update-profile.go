package usecase

import (
	"context"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/murilo-bracero/raspstore/commons/pkg/logger"
	"github.com/murilo-bracero/raspstore/idp/internal/model"
	"github.com/murilo-bracero/raspstore/idp/internal/repository"
)

type UpdateProfileUseCase interface {
	Execute(ctx context.Context, user *model.User) error
}

type updateProfileUseCase struct {
	userRepository repository.UsersRepository
}

func NewUpdateProfileUseCase(userRepository repository.UsersRepository) UpdateProfileUseCase {
	return &updateProfileUseCase{userRepository: userRepository}
}

func (u *updateProfileUseCase) Execute(ctx context.Context, user *model.User) error {
	traceId := ctx.Value(middleware.RequestIDKey).(string)

	found, err := u.userRepository.FindByUserId(user.UserId)

	if err != nil {
		logger.Error("[%s] Could not search for user: userId=%s: %s", traceId, user.UserId, err.Error())
		return err
	}

	found.Username = user.Username
	*user = *found

	if err := u.userRepository.Update(user); err != nil {
		logger.Error("[%s] Could not update user in database: userId=%s : %s",
			traceId, user.UserId, err.Error())
		return err
	}

	return nil
}
