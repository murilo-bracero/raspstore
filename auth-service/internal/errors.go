package internal

import "errors"

var (
	ErrUserNotFound         = errors.New("user with provided info does not exists")
	ErrIncorrectCredentials = errors.New("username or password does not match")
	ErrInvalidBasicAuth     = errors.New("invalid basic authorization header")
	ErrEmptyToken           = errors.New("token must not be empty")
	ErrConflict             = errors.New("username has already in use by another user")
	ErrInactiveAccount      = errors.New("account is inactive")
	ErrInvalidUsername      = errors.New("username should not be null or empty")
	ErrInvalidPassword      = errors.New("password should not be null or empty")
	ErrInvalidRoles         = errors.New("roles should not be null or empty")
)
