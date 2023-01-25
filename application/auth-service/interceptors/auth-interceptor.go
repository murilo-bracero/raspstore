package interceptor

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"raspstore.github.io/authentication/model"
	"raspstore.github.io/authentication/repository"
	"raspstore.github.io/authentication/token"
	"raspstore.github.io/authentication/utils"
)

var whitelistRoutes = "/pb.AuthService/login,/pb.AuthService/createCredentials,/pb.AuthService/authenticate"

type AuthInterceptor interface {
	WithUnaryAuthentication(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error)
	WithStreamingAuthentication(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error
}

type authInterceptor struct {
	tokenManager          token.TokenManager
	credentialsRepository repository.CredentialsRepository
}

func NewAuthInterceptor(tokenManager token.TokenManager, credentialsRepository repository.CredentialsRepository) AuthInterceptor {
	return &authInterceptor{tokenManager: tokenManager, credentialsRepository: credentialsRepository}
}

func (a *authInterceptor) WithStreamingAuthentication(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	if isRouteWhitelisted(info.FullMethod) {
		return handler(srv, ss)
	}

	if _, err := a.validateMetadata(ss.Context()); err != nil {
		return err
	}

	return handler(srv, ss)
}

func (a *authInterceptor) WithUnaryAuthentication(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {

	if isRouteWhitelisted(info.FullMethod) {
		return handler(ctx, req)
	}

	uid, err := a.validateMetadata(ctx)

	if err != nil {
		return nil, err
	}

	ctx = context.WithValue(ctx, utils.ContextKeyUID, uid)

	return handler(ctx, req)
}

func (a *authInterceptor) validateMetadata(ctx context.Context) (uid string, err error) {
	accessToken := utils.GetValueFromMetadata("authorization", ctx)

	if accessToken == "" {
		return "", status.Errorf(codes.Unauthenticated, "token invalid or malformed")
	}

	uid, err = a.tokenManager.Verify(accessToken)

	if err != nil {
		return "", status.Errorf(codes.Unauthenticated, "token invalid or malformed")
	}

	// just verifies if credentials exists on database, securing against fraud generated tokens

	var cred *model.Credential

	if cred, err = a.credentialsRepository.FindByUserId(uid); err != nil {
		return "", status.Errorf(codes.Unauthenticated, "token invalid or malformed")
	}

	fmt.Println("user ", cred.UserId, "authenticated")

	return uid, nil
}

func isRouteWhitelisted(route string) bool {
	for _, value := range strings.Split(whitelistRoutes, ",") {
		if value == route {
			return true
		}
	}

	return false
}
