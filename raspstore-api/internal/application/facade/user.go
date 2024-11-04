package facade

import (
	"log/slog"

	"github.com/google/uuid"
	"github.com/murilo-bracero/raspstore/file-service/internal/auth"
	"github.com/murilo-bracero/raspstore/file-service/internal/domain/entity"
	"github.com/murilo-bracero/raspstore/file-service/internal/domain/model"
	"github.com/murilo-bracero/raspstore/file-service/internal/infra/config"
	"github.com/murilo-bracero/raspstore/file-service/internal/infra/repository"
)

type UserFacade interface {
	Save(user *model.CreateUserRequest) error
}

type userFacade struct {
	config         *config.Config
	userRepository repository.Repository[entity.User]
}

func NewUserFacade(config *config.Config, userRepository repository.Repository[entity.User]) *userFacade {
	return &userFacade{
		config:         config,
		userRepository: userRepository,
	}
}

func (f *userFacade) Save(req *model.CreateUserRequest) error {
	hash, err := auth.HashPassword(req.Password)

	if err != nil {
		slog.Error("Could not hash password", "error", err)
		return err
	}

	rt, err := auth.GenerateRefreshToken(f.config)

	if err != nil {
		slog.Error("Could not generate refresh token", "error", err)
		return err
	}

	err = f.userRepository.Save(&entity.User{
		Id:           uuid.New(),
		Username:     req.Username,
		Name:         req.Name,
		PasswordHash: hash,
		RefreshToken: rt,
	})

	if err != nil {
		slog.Error("Could not save user", "error", err)
	}

	return err
}
