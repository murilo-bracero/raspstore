package internal

import "errors"

var (
	ErrEmptyUsername                = errors.New("field username could not be null or empty")
	ErrEmptyEmail                   = errors.New("field email could not be null or empty")
	ErrEmptyPassword                = errors.New("field password could not be null or empty")
	ErrComplexityPassword           = errors.New("field password must have at least %s characters")
	ErrEmailOrUsernameInUse         = errors.New("provided email or username has already in use by another user")
	ErrUserNotFound                 = errors.New("user with provided info does not exists")
	ErrInvalidTotpToken             = errors.New("2FA token is empty or invalid")
	ErrUserAlreadyExists            = errors.New("user with provided email or username already exists")
	ErrPublicUserCreationNotAllowed = errors.New("public user creation not allowed")
)
