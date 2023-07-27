package usecase

import (
	"context"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/murilo-bracero/raspstore/auth-service/internal"
	"github.com/murilo-bracero/raspstore/auth-service/internal/model"
	"github.com/murilo-bracero/raspstore/auth-service/internal/repository"
	"github.com/murilo-bracero/raspstore/commons/pkg/logger"
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

	if !found.IsEnabled {
		logger.Warn("[%s] User account is inactive: userId=%s", traceId, user.UserId)
		return internal.ErrInactiveAccount
	}

	if err != nil {
		logger.Error("[%s] Could not search for user: userId=%s: %s", traceId, user.UserId, err.Error())
		return err
	}

	if user.Username != found.Username {
		if exists, err := u.userRepository.ExistsByUsername(user.Username); err != nil {
			logger.Error("[%s] Could not check for conflicts in user new username: userId=%s, username=%s: %s",
				traceId, user.UserId, user.Username, err.Error())
			return err
		} else if exists {
			return internal.ErrConflict
		}
	}

	found.Username = user.Username
	*user = *found

	if err := u.userRepository.Update(user); err != nil {
		logger.Error("[%s] Could not update user in database: userId=%s, : %s",
			traceId, user.UserId, err.Error())
		return err
	}

	return nil
}
