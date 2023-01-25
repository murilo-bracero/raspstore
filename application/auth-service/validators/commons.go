package validators

import "errors"

var (
	ErrInvalidId            = errors.New("field Id must not be null or empty")
	ErrEmptyUsername        = errors.New("field username could not be null or empty")
	ErrEmptyEmail           = errors.New("field email could not be null or empty")
	ErrEmptyPassword        = errors.New("field password could not be null or empty")
	ErrComplexityPassword   = errors.New("field password must have at least 8 characters")
	ErrEmailOrUsernameInUse = errors.New("provided email or username has already in use by another user")
	ErrUserNotFound         = errors.New("user with provided info does not exists")
	ErrIncorrectCredentials = errors.New("provided email or password does not match")
	ErrEmptyToken           = errors.New("token must not be empty")
	ErrInvalidTotpToken     = errors.New("2FA token is empty or invalid")
	ErrEmptyUserId          = errors.New("userId could not be null or empty")
)

func GetErrorsList() []error {
	return []error{ErrInvalidId,
		ErrEmptyUsername,
		ErrEmptyEmail,
		ErrEmptyPassword,
		ErrComplexityPassword,
		ErrEmailOrUsernameInUse,
		ErrUserNotFound,
		ErrIncorrectCredentials,
		ErrEmptyToken}
}
