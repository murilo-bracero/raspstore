package validators

import (
	"raspstore.github.io/authentication/pb"
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
