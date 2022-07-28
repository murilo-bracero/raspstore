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

func ValidateCreateCredentials(req *pb.CreateCredentialsRequest) error {
	if req.Email == "" {
		return ErrEmptyEmail
	}

	if req.Hash == "" {
		return ErrEmptyPassword
	}

	if req.UserId == "" {
		return ErrEmptyUserId
	}

	return nil
}
