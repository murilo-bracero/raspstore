package service

import (
	"context"

	"github.com/murilo-bracero/raspstore-protofiles/users-service/pb"
	"raspstore.github.io/users-service/model"
	"raspstore.github.io/users-service/repository"
	"raspstore.github.io/users-service/validators"
)

type usersService struct {
	userRepository repository.UsersRepository
	pb.UnimplementedUsersServiceServer
}

func NewUserService(usersRepository repository.UsersRepository) pb.UsersServiceServer {
	return &usersService{userRepository: usersRepository}
}

func (u *usersService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.User, error) {
	err := validators.ValidateSignUp(req)

	if err != nil {
		return nil, err
	}

	found, err := u.userRepository.FindByEmail(req.Email)

	if err != nil {
		return nil, err
	}

	if found == nil {
		user := new(model.User)
		user.FromProtoBuffer(req)
		if err := u.userRepository.Save(user); err != nil {
			return nil, err
		}

		return user.ToProtoBuffer(), nil
	}

	return nil, validators.ErrEmailOrUsernameInUse
}

func (u *usersService) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error) {
	usr, err := u.userRepository.FindById(req.Id)

	if usr == nil {
		return nil, validators.ErrUserNotFound
	}

	return usr.ToProtoBuffer(), err
}

func (u *usersService) DeleteUser(ctx context.Context, req *pb.GetUserRequest) (*pb.DeleteUserResponse, error) {
	if err := u.userRepository.Delete(req.Id); err != nil {
		return nil, err
	}

	return &pb.DeleteUserResponse{Id: req.Id}, nil
}

func (u *usersService) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.User, error) {
	if err := validators.ValidateUpdate(req); err != nil {
		return nil, err
	}

	usr := new(model.User)
	if err := usr.FromUpdateProto(req); err != nil {
		return nil, err
	}

	if err := u.userRepository.Update(usr); err != nil {
		return nil, err
	}

	return usr.ToProtoBuffer(), nil
}

func (u *usersService) ListUser(req *pb.ListUsersRequest, stream pb.UsersService_ListUserServer) error {
	users, err := u.userRepository.FindAll()

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
