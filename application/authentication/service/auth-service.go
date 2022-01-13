package service

import (
	"context"
	"errors"
	"log"

	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"
	"raspstore.github.io/authentication/pb"
	"raspstore.github.io/authentication/repository"
	"raspstore.github.io/authentication/token"
	"raspstore.github.io/authentication/utils"
	"raspstore.github.io/authentication/validators"
)

type authService struct {
	userRepository repository.UsersRepository
	credRepository repository.CredentialsRepository
	tokenManager   token.TokenManager
	pb.UnimplementedAuthServiceServer
}

func NewAuthService(usersRepository repository.UsersRepository, credRepository repository.CredentialsRepository, tokenManager token.TokenManager) pb.AuthServiceServer {
	return &authService{userRepository: usersRepository, credRepository: credRepository, tokenManager: tokenManager}
}

func (a *authService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	if err := validators.ValidateLogin(req); err != nil {
		return nil, err
	}

	totpToken := ""

	if a.credRepository.Has2FAEnabledByEmail(req.Email) {
		totpToken = utils.GetValueFromMetadata("authorization", ctx)

	}

	if a.credRepository.IsCredentialsCorrect(req.Email, req.Password, totpToken) {
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

func (a *authService) Enroll2FA(ctx context.Context, req *pb.Enroll2FARequest) (*pb.Enroll2FAResponse, error) {
	uid := ctx.Value(utils.ContextKeyUID).(string)

	cred, err := a.credRepository.FindById(uid)

	if err != nil {
		return nil, err
	}

	if cred.Has2FAEnabled {
		cred.Has2FAEnabled = false
		cred.Secret = ""
		if err := a.credRepository.Update(cred); err != nil {
			log.Println("could not update 2FA credentials for user ", uid, " error: ", err.Error())
			return nil, err
		}
	}

	key, err := totp.Generate(totp.GenerateOpts{Issuer: "raspstore", AccountName: cred.Email})

	if err != nil {
		log.Println("could not generate user totp key for user ", uid, " error: ", err.Error())
		return nil, err
	}

	return &pb.Enroll2FAResponse{Secret: key.Secret()}, nil
}

func (a *authService) Verify2FA(ctx context.Context, req *pb.Verify2FARequest) (*pb.Verify2FAResponse, error) {
	if err := validators.Validate2FARequest(req); err != nil {
		return nil, err
	}

	uid := ctx.Value(utils.ContextKeyUID).(string)

	cred, err := a.credRepository.FindById(uid)

	if err != nil {
		return nil, err
	}

	if cred.Has2FAEnabled {
		return nil, errors.New("user already have 2FA enrolled and in use. For authenticate an user with 2FA enables, use Login method")
	}

	isValid := totp.Validate(req.Token, cred.Secret)

	return &pb.Verify2FAResponse{Status: isValid}, nil
}

func hash(text string) (hash string, err error) {
	bts := []byte(text)

	raw, err := bcrypt.GenerateFromPassword(bts, bcrypt.DefaultCost)
	return string(raw[:]), err
}
