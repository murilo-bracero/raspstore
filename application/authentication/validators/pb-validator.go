package validators

import (
	"errors"

	"raspstore.github.io/authentication/pb"
)

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
)

func ValidateLogin(req *pb.LoginRequest) error {
	if req.Email == "" {
		return ErrEmptyEmail
	}

	if req.Password == "" {
		return ErrEmptyPassword
	}

	return nil
}

func ValidateAuthenticate(req *pb.AuthenticateRequest) error {
	if req.Token == "" {
		return ErrEmptyToken
	}
	return nil
}

func ValidateSignUp(req *pb.CreateUserRequest) error {
	if req.Email == "" {
		return ErrEmptyEmail
	}

	if req.Username == "" {
		return ErrEmptyUsername
	}

	if req.Password == "" {
		return ErrEmptyPassword
	}

	if len(req.Password) < 8 {
		return ErrComplexityPassword
	}

	return nil
}

func ValidateUpdate(req *pb.UpdateUserRequest) error {
	if req.Id == "" {
		return ErrInvalidId
	}

	if req.Email == "" {
		return ErrEmptyEmail
	}

	if req.Username == "" {
		return ErrEmptyUsername
	}

	return nil
}
