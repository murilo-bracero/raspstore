package service

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"raspstore.github.io/authentication/model"
	"raspstore.github.io/authentication/pb"
	"raspstore.github.io/authentication/repository"
	"raspstore.github.io/authentication/validators"
)

type AuthService interface {
	SignUp(ctx context.Context, req *pb.CreateUserRequest) (*pb.User, error)
	GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error)
	DeleteUser(ctx context.Context, req *pb.GetUserRequest) (*pb.DeleteUserResponse, error)
	UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.User, error)
	ListUser(req *pb.ListUsersRequest, stream pb.AuthService_ListUserServer) error
}

type authService struct {
	userRepository repository.UsersRepository
	pb.UnimplementedAuthServiceServer
}

func NewAuthService(usersRepository repository.UsersRepository) pb.AuthServiceServer {
	return &authService{userRepository: usersRepository}
}

func (a *authService) SignUp(ctx context.Context, req *pb.CreateUserRequest) (*pb.User, error) {
	err := validators.ValidateSignUp(req)

	if err != nil {
		return nil, err
	}

	found, err := a.userRepository.FindByEmailOrUsername(req.Email, req.Username)

	if err == mongo.ErrNoDocuments {
		user := new(model.User)
		user.FromProtoBuffer(req)
		err := a.userRepository.Save(user)

		if err != nil {
			return nil, err
		}

		return user.ToProtoBuffer(), nil
	}

	if found == nil {
		return nil, err
	}

	return nil, validators.ErrEmailOrUsernameInUse
}

func (a *authService) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error) {
	usr, err := a.userRepository.FindById(req.Id)
	return usr.ToProtoBuffer(), err
}

func (a *authService) DeleteUser(ctx context.Context, req *pb.GetUserRequest) (*pb.DeleteUserResponse, error) {
	err := a.userRepository.DeleteUser(req.Id)

	if err != nil {
		return nil, err
	}

	return &pb.DeleteUserResponse{Id: req.Id}, nil
}

func (a *authService) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.User, error) {
	err := validators.ValidateUpdate(req)

	if err != nil {
		return nil, err
	}

	usr := new(model.User)
	err = usr.FromUpdateProto(req)

	if err != nil {
		return nil, err
	}

	err = a.userRepository.UpdateUser(usr)

	if err != nil {
		return nil, err
	}

	var found *model.User
	found, err = a.userRepository.FindById(req.Id)

	if err != nil {
		return nil, err
	}

	return found.ToProtoBuffer(), nil
}
func (a *authService) ListUser(req *pb.ListUsersRequest, stream pb.AuthService_ListUserServer) error {
	users, err := a.userRepository.FindAll()

	if err != nil {
		return err
	}

	for _, user := range users {
		err := stream.Send(user.ToProtoBuffer())
		if err != nil {
			return err
		}
	}

	return nil
}
