package interceptor

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"raspstore.github.io/authentication/repository"
	"raspstore.github.io/authentication/token"
)

var whitelistRoutes = "/pb.AuthService/Login,/pb.AuthService/SignUp,/pb.AuthService/Authenticate"

type AuthInterceptor interface {
	WithUnaryAuthentication(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error)
	WithStreamingAuthentication(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error
}

type authInterceptor struct {
	repo         repository.UsersRepository
	tokenManager token.TokenManager
}

func NewAuthInterceptor(repo repository.UsersRepository, tokenManager token.TokenManager) AuthInterceptor {
	return &authInterceptor{repo: repo, tokenManager: tokenManager}
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

	if _, err := a.validateMetadata(ctx); err != nil {
		return nil, err
	}

	return handler(ctx, req)
}

func (a *authInterceptor) validateMetadata(ctx context.Context) (uid string, err error) {
	md, exists := metadata.FromIncomingContext(ctx)

	if !exists {
		return "", status.Errorf(codes.Unauthenticated, "metadata not provided")
	}

	values := md["authorization"]

	if len(values) == 0 {
		return "", status.Errorf(codes.Unauthenticated, "authorization metadata not provided")
	}

	accessToken := values[0]

	uid, err = a.tokenManager.Verify(accessToken)

	if err != nil {
		return "", status.Errorf(codes.Unauthenticated, "token invalid or expired")
	}

	// just verifies if user exists on database, securing against fraud generated tokens

	if usr, err := a.repo.FindById(uid); err != nil {
		return "", status.Errorf(codes.Unauthenticated, "fraudulent token")
	} else {
		fmt.Println("user ", usr.Email, "authenticated")
	}

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
