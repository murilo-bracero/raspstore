package grpc

import (
	"context"
	"strings"

	"github.com/murilo-bracero/raspstore/auth-service/internal"
	"github.com/murilo-bracero/raspstore/auth-service/internal/token"
	"github.com/murilo-bracero/raspstore/auth-service/proto/v1/auth-service/pb"
)

const tokenScheme = "Bearer "

type authService struct {
	pb.UnimplementedAuthServiceServer
}

func NewAuthService() pb.AuthServiceServer {
	return &authService{}
}

func (a *authService) Authenticate(ctx context.Context, req *pb.AuthenticateRequest) (*pb.AuthenticateResponse, error) {
	if err := validateAuthenticateRequest(req); err != nil {
		return nil, err
	}

	if claims, err := token.Verify(strings.ReplaceAll(req.Token, tokenScheme, "")); err != nil {
		return nil, err
	} else {
		return &pb.AuthenticateResponse{Uid: claims.Subject, Roles: claims.Roles}, nil
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
