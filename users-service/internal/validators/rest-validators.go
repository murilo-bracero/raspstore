package validators

import (
	v1 "github.com/murilo-bracero/raspstore/users-service/api/v1"
	"github.com/murilo-bracero/raspstore/users-service/internal"
)

func ValidateCreateUserRequest(req v1.CreateUserRequest) error {
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
