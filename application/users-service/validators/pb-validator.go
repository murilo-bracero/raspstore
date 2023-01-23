package validators

import "github.com/murilo-bracero/raspstore-protofiles/users-service/pb"

func ValidateSignUp(req *pb.CreateUserRequest) error {
	if req.Email == "" {
		return ErrEmptyEmail
	}

	if req.Username == "" {
		return ErrEmptyUsername
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
