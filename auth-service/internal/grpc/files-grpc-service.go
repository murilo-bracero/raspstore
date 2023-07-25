package grpc

import (
	"context"
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/murilo-bracero/raspstore/auth-service/internal"
	"github.com/murilo-bracero/raspstore/auth-service/internal/model"
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

	if claims, err := verifyToken(strings.ReplaceAll(req.Token, tokenScheme, "")); err != nil {
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

func verifyToken(rawToken string) (userClaims *model.UserClaims, err error) {

	parsedToken, err := jwt.ParseWithClaims(rawToken, &model.UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("error reading jwt: wrong signing method: %v", token.Header["alg"])
		}
		return []byte(internal.TokenSecret()), nil
	})

	if err != nil {
		return nil, err
	}

	return parsedToken.Claims.(*model.UserClaims), nil
}
