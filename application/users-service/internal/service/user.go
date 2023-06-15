package service

import (
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"raspstore.github.io/users-service/internal"
	"raspstore.github.io/users-service/internal/model"
	"raspstore.github.io/users-service/internal/repository"
)

type UserService interface {
	CreateUser(user *model.User) error
	GetUserById(id string) (*model.User, error)
	GetAllUsersByPage(page int, size int) (*model.UserPage, error)
	UpdateUser(user *model.User) (*model.User, error)
	RemoveUserById(id string) error
}

type userService struct {
	usersRepository  repository.UsersRepository
	configRepository repository.UsersConfigRepository
}

func NewUserService(usersRepository repository.UsersRepository, configRepository repository.UsersConfigRepository) UserService {
	return &userService{usersRepository: usersRepository, configRepository: configRepository}
}

func (s *userService) CreateUser(user *model.User) error {

	if exists, err := s.usersRepository.ExistsByEmailOrUsername(user.Email, user.Username); err != nil {
		return err
	} else if exists {
		return internal.ErrUserAlreadyExists
	}

	if hash, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost); err != nil {
		return err
	} else {
		user.PasswordHash = string(hash)
	}

	if err := s.usersRepository.Save(user); err != nil {
		return err
	}

	return nil
}

func (s *userService) GetUserById(id string) (*model.User, error) {
	user, err := s.usersRepository.FindById(id)

	if err == mongo.ErrNoDocuments {
		return nil, internal.ErrUserNotFound
	}

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) GetAllUsersByPage(page int, size int) (*model.UserPage, error) {
	userPage, err := s.usersRepository.FindAll(page, size)

	if err != nil {
		return nil, err
	}

	return userPage, nil
}

func (s *userService) UpdateUser(user *model.User) (*model.User, error) {
	found, err := s.usersRepository.FindById(user.UserId)

	if err == mongo.ErrNoDocuments {
		return nil, internal.ErrUserNotFound
	}

	if err != nil {
		return nil, err
	}

	if user.Username != found.Username || user.Email != found.Email {
		if exists, err := s.usersRepository.ExistsByEmailOrUsername(user.Email, user.Username); err != nil {
			return nil, err
		} else if exists {
			return nil, internal.ErrEmailOrUsernameInUse
		}
	}

	if err := s.usersRepository.Update(user); err != nil {
		return nil, err
	}

	found.Username = user.Username
	found.Email = user.Email
	found.PhoneNumber = user.PhoneNumber

	return found, nil
}

func (s *userService) RemoveUserById(id string) error {
	return s.usersRepository.Delete(id)
}
