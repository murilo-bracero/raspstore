package internal

import "errors"

var (
	ErrUserNotFound         = errors.New("user with provided info does not exists")
	ErrIncorrectCredentials = errors.New("username or password does not match")
	ErrInvalidBasicAuth     = errors.New("invalid basic authorization header")
	ErrEmptyToken           = errors.New("token must not be empty")
	ErrConflict             = errors.New("username has already in use by another user")
)
