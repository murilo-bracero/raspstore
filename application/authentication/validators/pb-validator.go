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
	ErrEmailOrUsernameInUse = errors.New("provided email or username has already in use by another user")
)

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
