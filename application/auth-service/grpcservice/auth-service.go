package service

import (
	"context"

	"github.com/murilo-bracero/raspstore-protofiles/auth-service/pb"
	"raspstore.github.io/auth-service/token"
	"raspstore.github.io/auth-service/utils"
)

type authService struct {
	tokenManager token.TokenManager
	pb.UnimplementedAuthServiceServer
}

func NewAuthService(tokenManager token.TokenManager) pb.AuthServiceServer {
	return &authService{tokenManager: tokenManager}
}

func (a *authService) Authenticate(ctx context.Context, req *pb.AuthenticateRequest) (*pb.AuthenticateResponse, error) {
	if err := validateAuthenticateRequest(req); err != nil {
		return nil, err
	}

	if uid, err := a.tokenManager.Verify(req.Token); err != nil {
		return nil, err
	} else {
		return &pb.AuthenticateResponse{Uid: uid}, nil
	}
}

func validateAuthenticateRequest(req *pb.AuthenticateRequest) error {
	if req.Token == "" {
		return utils.ErrEmptyToken
	}
	return nil
}
