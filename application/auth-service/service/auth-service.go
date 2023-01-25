package service

import (
	"context"
	"log"

	"github.com/murilo-bracero/raspstore-protofiles/authentication/pb"
	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"
	"raspstore.github.io/authentication/model"
	"raspstore.github.io/authentication/repository"
	"raspstore.github.io/authentication/token"
	"raspstore.github.io/authentication/utils"
	"raspstore.github.io/authentication/validators"
)

type authService struct {
	credRepository repository.CredentialsRepository
	tokenManager   token.TokenManager
	pb.UnimplementedAuthServiceServer
}

func NewAuthService(credRepository repository.CredentialsRepository, tokenManager token.TokenManager) pb.AuthServiceServer {
	return &authService{credRepository: credRepository, tokenManager: tokenManager}
}

func (a *authService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	if err := validators.ValidateLogin(req); err != nil {
		return nil, err
	}

	cred, err := a.credRepository.FindByEmail(req.Email)

	if err != nil {
		return nil, validators.ErrIncorrectCredentials
	}

	if !isValidPassword(req.Password, cred.Hash) {
		return nil, validators.ErrIncorrectCredentials
	}

	if cred.Has2FAEnabled {
		totpToken := utils.GetValueFromMetadata("authorization", ctx)

		if !isValidTotp(totpToken, cred.Secret) {
			return nil, validators.ErrIncorrectCredentials
		}
	}

	if token, err := a.tokenManager.Generate(cred.UserId); err != nil {
		log.Println("error while generating token for user ", cred.UserId, ":", err)
		return nil, validators.ErrIncorrectCredentials
	} else {
		return &pb.LoginResponse{Token: token}, nil
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

func (a *authService) CreateCredentials(ctx context.Context, req *pb.CreateCredentialsRequest) (*pb.CreateCredentialsResponse, error) {
	if err := validators.ValidateCreateCredentials(req); err != nil {
		return nil, err
	}

	cred := model.ConvertToModel(req)

	if err := a.credRepository.Save(cred); err != nil {
		return nil, err
	}

	return &pb.CreateCredentialsResponse{
		CredentialsId: cred.Id.Hex(),
	}, nil
}

func isValidPassword(rawPass string, hashPass string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashPass), []byte(rawPass)) == nil
}

func isValidTotp(token string, secret string) bool {
	return totp.Validate(token, secret)
}
