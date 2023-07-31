package usecase

import (
	"context"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/murilo-bracero/raspstore/commons/pkg/logger"
	"github.com/murilo-bracero/raspstore/idp/internal/model"
	"github.com/murilo-bracero/raspstore/idp/internal/repository"
)

type UpdateUserUseCase interface {
	Execute(ctx context.Context, user *model.User) error
}

type updateUserUseCase struct {
	userRepository repository.UsersRepository
}

func NewUpdateUserUseCase(userRepository repository.UsersRepository) UpdateUserUseCase {
	return &updateUserUseCase{userRepository: userRepository}
}

func (u *updateUserUseCase) Execute(ctx context.Context, user *model.User) error {
	traceId := ctx.Value(middleware.RequestIDKey).(string)

	found, err := u.userRepository.FindByUserId(user.UserId)

	if err != nil {
		logger.Error("[%s] Could not search for user: userId=%s: %s", traceId, user.UserId, err.Error())
		return err
	}

	found.Username = user.Username
	found.IsEnabled = user.IsEnabled
	found.Permissions = user.Permissions

	if !user.IsMfaEnabled {
		found.IsMfaEnabled = false
		found.IsMfaVerified = false
		found.Secret = ""
		found.RefreshToken = ""
	}

	*user = *found

	if err := u.userRepository.Update(user); err != nil {
		logger.Error("[%s] Could not update user in database: userId=%s : %s",
			traceId, user.UserId, err.Error())
		return err
	}

	return nil
}
