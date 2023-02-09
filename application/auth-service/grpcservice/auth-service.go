package service

import (
	"context"
	"strings"

	"github.com/murilo-bracero/raspstore-protofiles/auth-service/pb"
	"raspstore.github.io/auth-service/token"
	"raspstore.github.io/auth-service/utils"
)

const tokenScheme = "Bearer"

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

	if uid, err := a.tokenManager.Verify(strings.ReplaceAll(req.Token, "Bearer ", "")); err != nil {
		return nil, err
	} else {
		return &pb.AuthenticateResponse{Uid: uid}, nil
	}
}

func validateAuthenticateRequest(req *pb.AuthenticateRequest) error {
	if req.Token == "" {
		return utils.ErrEmptyToken
	}

	tokenParts := strings.Split(req.Token, " ")

	if len(tokenParts) != 2 {
		return utils.ErrIncorrectCredentials
	}

	if tokenParts[0] != tokenScheme {
		return utils.ErrIncorrectCredentials
	}

	return nil
}
