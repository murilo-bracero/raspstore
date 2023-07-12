package service

import (
	"fmt"

	"raspstore.github.io/users-service/internal"
	"raspstore.github.io/users-service/internal/model"
	"raspstore.github.io/users-service/internal/repository"
)

type UserConfigService interface {
	UpdateUserConfig(userConfig *model.UserConfiguration) error
	GetUserConfig() (userConfig *model.UserConfiguration, err error)
	ValidateUser(user *model.User, isAdminCreation bool) error
}

type userConfigService struct {
	configRepository repository.UsersConfigRepository
}

func NewUserConfigService(configRepository repository.UsersConfigRepository) UserConfigService {
	return &userConfigService{configRepository: configRepository}
}

func (s *userConfigService) UpdateUserConfig(userConfig *model.UserConfiguration) error {
	return s.configRepository.Update(userConfig)
}

func (s *userConfigService) GetUserConfig() (userConfig *model.UserConfiguration, err error) {
	return s.configRepository.Find()
}

func (s *userConfigService) ValidateUser(user *model.User, isAdminCreation bool) error {
	userConf, err := s.GetUserConfig()

	if err != nil {
		return err
	}

	if !isAdminCreation && userConf.AllowPublicUserCreation {
		return internal.ErrPublicUserCreationNotAllowed
	}

	if len(user.PasswordHash) < userConf.MinPasswordLength {
		return fmt.Errorf(internal.ErrComplexityPassword.Error(), userConf.MinPasswordLength)
	}

	return nil
}
