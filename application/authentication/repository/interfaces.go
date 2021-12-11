package repository

import "raspstore.github.io/authentication/model"

type UsersRepository interface {
	Save(user *model.User) error
	FindById(id string) (user *model.User, err error)
	FindByEmail(email string) (user *model.User, err error)
	DeleteUser(id string) error
	UpdateUser(user *model.User) error
	FindAll() (users []*model.User, err error)
}

type CredentialsRepository interface {
	Save(user *model.User, password string) error
	Update(user *model.User) error
	Delete(id string) error
	IsCredentialsCorrect(email string, password string) bool
}
