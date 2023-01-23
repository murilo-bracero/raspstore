package interceptor

import (
	"context"
	"log"
	"strings"

	"github.com/murilo-bracero/raspstore-protofiles/authentication/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"raspstore.github.io/users-service/db"
)

var whitelistRoutes = "/pb.UsersService/SignUp"

type UidKey string

type AuthInterceptor interface {
	WithUnaryAuthentication(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error)
	WithStreamingAuthentication(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error
}

type authInterceptor struct {
	cfg db.Config
}

func NewAuthInterceptor(cfg db.Config) AuthInterceptor {
	return &authInterceptor{cfg: cfg}
}

func (a *authInterceptor) WithStreamingAuthentication(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	if isRouteWhitelisted(info.FullMethod) {
		return handler(srv, ss)
	}

	accessToken, err := getTokenFromContext(ss.Context())

	if err != nil {
		return err
	}

	uid, err := a.validateToken(accessToken)

	if err != nil {
		return err
	}

	log.Println("user ", uid, " streammed resource ", info.FullMethod)

	return handler(srv, ss)
}

func (a *authInterceptor) WithUnaryAuthentication(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {

	if isRouteWhitelisted(info.FullMethod) {
		return handler(ctx, req)
	}

	accessToken, err := getTokenFromContext(ctx)

	if err != nil {
		return nil, err
	}

	uid, err := a.validateToken(accessToken)

	if err != nil {
		return nil, err
	}

	log.Println("user ", uid, " accessed resource ", info.FullMethod)

	ctx = context.WithValue(ctx, UidKey("uid"), uid)

	return handler(ctx, req)
}

func (a *authInterceptor) validateToken(token string) (uid string, err error) {
	conn, err := grpc.Dial(a.cfg.AuthServiceUrl())

	if err != nil {
		log.Fatalln("could not stablish connection to auth service, it may goes down: ", err.Error())
	}

	defer conn.Close()

	client := pb.NewAuthServiceClient(conn)

	in := &pb.AuthenticateRequest{Token: token}

	if res, err := client.Authenticate(context.Background(), in); err != nil {
		return "", err
	} else {
		return res.Uid, nil
	}
}

func getTokenFromContext(ctx context.Context) (string, error) {
	md, exists := metadata.FromIncomingContext(ctx)

	if !exists {
		return "", status.Errorf(codes.Unauthenticated, "metadata not provided")
	}

	values := md["authorization"]

	if len(values) == 0 {
		return "", status.Errorf(codes.Unauthenticated, "authorization metadata not provided")
	}

	return values[0], nil
}

func isRouteWhitelisted(route string) bool {
	for _, value := range strings.Split(whitelistRoutes, ",") {
		if value == route {
			return true
		}
	}

	return false
}
