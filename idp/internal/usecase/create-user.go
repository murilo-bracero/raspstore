package usecase

import (
	"context"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/murilo-bracero/raspstore/commons/pkg/logger"
	"github.com/murilo-bracero/raspstore/idp/internal/infra"
	"github.com/murilo-bracero/raspstore/idp/internal/model"
	"github.com/murilo-bracero/raspstore/idp/internal/repository"
)

type CreateUserUseCase interface {
	Execute(ctx context.Context, user *model.User) error
}

type createUserUseCase struct {
	config         *infra.Config
	userRepository repository.UsersRepository
}

func NewCreateUserUseCase(userRepository repository.UsersRepository, config *infra.Config) CreateUserUseCase {
	return &createUserUseCase{userRepository: userRepository, config: config}
}

func (u *createUserUseCase) Execute(ctx context.Context, user *model.User) error {
	traceId := ctx.Value(middleware.RequestIDKey).(string)

	if u.config.EnforceMfa {
		user.IsMfaEnabled = true
	}
	err := u.userRepository.Save(user)

	if err != nil {
		logger.Error("[%s] Could create user: %s", traceId, err.Error())
		return err
	}

	return nil
}
