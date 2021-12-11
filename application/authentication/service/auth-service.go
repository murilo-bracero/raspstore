package service

import (
	"context"
	"fmt"

	"raspstore.github.io/authentication/model"
	"raspstore.github.io/authentication/pb"
	"raspstore.github.io/authentication/repository"
	"raspstore.github.io/authentication/token"
	"raspstore.github.io/authentication/validators"
)

type AuthService interface {
	SignUp(ctx context.Context, req *pb.CreateUserRequest) (*pb.User, error)
	GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error)
	DeleteUser(ctx context.Context, req *pb.GetUserRequest) (*pb.DeleteUserResponse, error)
	UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.User, error)
	ListUser(req *pb.ListUsersRequest, stream pb.AuthService_ListUserServer) error
	Authenticate(ctx context.Context, req *pb.AuthenticateRequest) (*pb.AuthenticateResponse, error)
	Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error)
}

type authService struct {
	userRepository repository.UsersRepository
	credRepository repository.CredentialsRepository
	tokenManager   token.TokenManager
	pb.UnimplementedAuthServiceServer
}

func NewAuthService(usersRepository repository.UsersRepository, credRepository repository.CredentialsRepository, tokenManager token.TokenManager) pb.AuthServiceServer {
	return &authService{userRepository: usersRepository, credRepository: credRepository, tokenManager: tokenManager}
}

func (a *authService) SignUp(ctx context.Context, req *pb.CreateUserRequest) (*pb.User, error) {
	err := validators.ValidateSignUp(req)

	if err != nil {
		return nil, err
	}

	found, err := a.userRepository.FindByEmail(req.Email)

	if err != nil {
		return nil, err
	}

	if found == nil {
		user := new(model.User)
		user.FromProtoBuffer(req)
		if err := a.userRepository.Save(user); err != nil {
			return nil, err
		}

		if err := a.credRepository.Save(user, req.Password); err != nil {
			if inner_error := a.userRepository.DeleteUser(user.UserId); inner_error != nil {
				fmt.Println("user of id ", user.UserId, " failed to be inserted in credentials database but was inserted in user database. Remove it manually")
				return nil, inner_error
			}
			return nil, err
		}

		return user.ToProtoBuffer(), nil
	}

	return nil, validators.ErrEmailOrUsernameInUse
}

func (a *authService) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error) {
	usr, err := a.userRepository.FindById(req.Id)

	if usr == nil {
		return nil, validators.ErrUserNotFound
	}

	return usr.ToProtoBuffer(), err
}

func (a *authService) DeleteUser(ctx context.Context, req *pb.GetUserRequest) (*pb.DeleteUserResponse, error) {
	if err := a.userRepository.DeleteUser(req.Id); err != nil {
		return nil, err
	}

	if err := a.credRepository.Delete(req.Id); err != nil {
		fmt.Print("User ", req.Id, " was removed from user database but not for credentials database, remove it mannually")
		return nil, err
	}

	return &pb.DeleteUserResponse{Id: req.Id}, nil
}

func (a *authService) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.User, error) {
	if err := validators.ValidateUpdate(req); err != nil {
		return nil, err
	}

	usr := new(model.User)
	if err := usr.FromUpdateProto(req); err != nil {
		return nil, err
	}

	if err := a.userRepository.UpdateUser(usr); err != nil {
		return nil, err
	}

	if err := a.credRepository.Update(usr); err != nil {
		fmt.Println("user of id ", usr.UserId, " failed to be updated in credentials database but was inserted in user database. Update it manually")
		return nil, err
	}

	if found, err := a.userRepository.FindById(req.Id); err != nil {
		return nil, err
	} else {
		return found.ToProtoBuffer(), nil
	}
}

func (a *authService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	if err := validators.ValidateLogin(req); err != nil {
		return nil, err
	}

	if a.credRepository.IsCredentialsCorrect(req.Email, req.Password) {
		user, err := a.userRepository.FindByEmail(req.Email)

		if err != nil {
			return nil, err
		}

		if user == nil {
			return nil, validators.ErrUserNotFound
		}

		if token, err := a.tokenManager.Generate(user.UserId); err != nil {
			return nil, err
		} else {
			return &pb.LoginResponse{Token: token}, nil
		}
	} else {
		return nil, validators.ErrIncorrectCredentials
	}

}

func (a *authService) Authenticate(ctx context.Context, req *pb.AuthenticateRequest) (*pb.AuthenticateResponse, error) {
	if err := validators.ValidateAuthenticate(req); err != nil {
		return nil, err
	}

	if uid, err := a.tokenManager.Verify(req.Token); err != nil {
		return nil, err
	} else {
		return &pb.AuthenticateResponse{Uid: uid}, nil
	}
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
