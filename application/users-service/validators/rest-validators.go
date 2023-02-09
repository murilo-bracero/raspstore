package validators

import (
	"raspstore.github.io/users-service/api/dto"
	"raspstore.github.io/users-service/internal"
)

func ValidateCreateUserRequest(req dto.CreateUserRequest) error {
	if req.Username == "" {
		return internal.ErrEmptyUsername
	}

	if req.Email == "" {
		return internal.ErrEmptyEmail
	}

	if req.Password == "" {
		return internal.ErrEmptyPassword
	}

	return nil
}
