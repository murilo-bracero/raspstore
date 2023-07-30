package usecase

import (
	"context"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/murilo-bracero/raspstore/auth-service/internal/model"
	"github.com/murilo-bracero/raspstore/auth-service/internal/repository"
	"github.com/murilo-bracero/raspstore/commons/pkg/logger"
	"golang.org/x/crypto/bcrypt"
)

type CreateUserUseCase interface {
	Execute(ctx context.Context, user *model.User) error
}

type createUserUseCase struct {
	userRepository repository.UsersRepository
}

func NewCreateUserUseCase(userRepository repository.UsersRepository) CreateUserUseCase {
	return &createUserUseCase{userRepository: userRepository}
}

func (u *createUserUseCase) Execute(ctx context.Context, user *model.User) error {
	traceId := ctx.Value(middleware.RequestIDKey).(string)

	if hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost); err != nil {
		logger.Error("[%s] Could not hash user password: %s", traceId, err.Error())
		return err
	} else {
		user.Password = string(hash)
	}

	err := u.userRepository.Save(user)

	if err != nil {
		logger.Error("[%s] Could create user: %s", traceId, err.Error())
		return err
	}

	return nil
}
