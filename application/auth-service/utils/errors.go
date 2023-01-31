package utils

import "errors"

var (
	ErrUserNotFound         = errors.New("user with provided info does not exists")
	ErrIncorrectCredentials = errors.New("provided email or password does not match")
	ErrEmptyToken           = errors.New("token must not be empty")
)

func GetErrorsList() []error {
	return []error{ErrUserNotFound, ErrIncorrectCredentials, ErrEmptyToken}
}
