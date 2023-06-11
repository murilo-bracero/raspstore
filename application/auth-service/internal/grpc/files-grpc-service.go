package grpc

import (
	"context"
	"strings"

	"github.com/murilo-bracero/raspstore-protofiles/auth-service/pb"
	"raspstore.github.io/auth-service/internal"
	"raspstore.github.io/auth-service/internal/service"
)

const tokenScheme = "Bearer"

type authService struct {
	tokenService service.TokenService
	pb.UnimplementedAuthServiceServer
}

func NewAuthService(tokenService service.TokenService) pb.AuthServiceServer {
	return &authService{tokenService: tokenService}
}

func (a *authService) Authenticate(ctx context.Context, req *pb.AuthenticateRequest) (*pb.AuthenticateResponse, error) {
	if err := validateAuthenticateRequest(req); err != nil {
		return nil, err
	}

	if uid, err := a.tokenService.Verify(strings.ReplaceAll(req.Token, "Bearer ", "")); err != nil {
		return nil, err
	} else {
		return &pb.AuthenticateResponse{Uid: uid}, nil
	}
}

func validateAuthenticateRequest(req *pb.AuthenticateRequest) error {
	if req.Token == "" {
		return internal.ErrEmptyToken
	}

	tokenParts := strings.Split(req.Token, " ")

	if len(tokenParts) != 2 {
		return internal.ErrIncorrectCredentials
	}

	if tokenParts[0] != tokenScheme {
		return internal.ErrIncorrectCredentials
	}

	return nil
}
