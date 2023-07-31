package grpc

import (
	"context"
	"strings"

	"github.com/murilo-bracero/raspstore/idp/internal"
	"github.com/murilo-bracero/raspstore/idp/internal/infra"
	"github.com/murilo-bracero/raspstore/idp/internal/token"
	"github.com/murilo-bracero/raspstore/idp/proto/v1/auth-service/pb"
)

const tokenScheme = "Bearer"

type authService struct {
	config *infra.Config
	pb.UnimplementedAuthServiceServer
}

func NewAuthService(config *infra.Config) pb.AuthServiceServer {
	return &authService{config: config}
}

func (a *authService) Authenticate(ctx context.Context, req *pb.AuthenticateRequest) (*pb.AuthenticateResponse, error) {
	if err := validateAuthenticateRequest(req); err != nil {
		return nil, err
	}

	claims, err := token.Verify(a.config, strings.ReplaceAll(req.Token, tokenScheme+" ", ""))

	if err != nil {
		return nil, err
	}

	return &pb.AuthenticateResponse{Uid: claims.Uid, Roles: claims.Roles}, nil
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
