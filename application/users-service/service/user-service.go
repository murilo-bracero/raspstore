package service

import (
	"context"
	"fmt"

	"raspstore.github.io/users-service/model"
	"raspstore.github.io/users-service/pb"
	"raspstore.github.io/users-service/repository"
	"raspstore.github.io/users-service/utils"
	"raspstore.github.io/users-service/validators"
)

type usersService struct {
	userRepository repository.UsersRepository
	credRepository repository.CredentialsRepository
	pb.UnimplementedUsersServiceServer
}

func NewUserService(usersRepository repository.UsersRepository, credRepository repository.CredentialsRepository) pb.UsersServiceServer {
	return &usersService{userRepository: usersRepository, credRepository: credRepository}
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

		hash, err := utils.Hash(req.Password)

		if err != nil {
			return nil, err
		}

		cred := &model.Credential{
			Id:            user.UserId,
			Email:         user.Email,
			Hash:          hash,
			Has2FAEnabled: false,
		}

		if err := u.credRepository.Save(cred); err != nil {
			if inner_error := u.userRepository.DeleteUser(user.UserId); inner_error != nil {
				fmt.Println("user of id ", user.UserId, " failed to be inserted in credentials database but was inserted in user database. Remove it manually")
				return nil, inner_error
			}
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
	if err := u.userRepository.DeleteUser(req.Id); err != nil {
		return nil, err
	}

	if err := u.credRepository.Delete(req.Id); err != nil {
		fmt.Print("User ", req.Id, " was removed from user database but not for credentials database, remove it mannually")
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

	if err := u.userRepository.UpdateUser(usr); err != nil {
		return nil, err
	}

	cred := &model.Credential{
		Id:    usr.UserId,
		Email: usr.Email,
	}

	if err := u.credRepository.Update(cred); err != nil {
		fmt.Println("user of id ", usr.UserId, " failed to be updated in credentials database but was inserted in user database. Update it manually")
		return nil, err
	}

	if found, err := u.userRepository.FindById(req.Id); err != nil {
		return nil, err
	} else {
		return found.ToProtoBuffer(), nil
	}
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
