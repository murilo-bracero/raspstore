package interceptor

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"raspstore.github.io/file-manager/db"
	"raspstore.github.io/file-manager/pb"
)

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
	md, exists := metadata.FromIncomingContext(ss.Context())

	if !exists {
		return status.Errorf(codes.Unauthenticated, "metadata not provided")
	}

	values := md["authorization"]

	if len(values) == 0 {
		return status.Errorf(codes.Unauthenticated, "authorization metadata not provided")
	}

	accessToken := values[0]

	conn, err := grpc.Dial(a.cfg.AuthServiceUrl(), grpc.WithInsecure())

	if err != nil {
		log.Fatalln("could not stablish connection to auth service, it may goes down: ", err.Error())
	}

	client := pb.NewAuthServiceClient(conn)

	in := &pb.AuthenticateRequest{Token: accessToken}

	if res, err := client.Authenticate(context.Background(), in); err != nil {
		return err
	} else {
		log.Println("user ", res.Uid, " accessed resource ", info.FullMethod)
	}

	defer conn.Close()

	return handler(srv, ss)
}

func (a *authInterceptor) WithUnaryAuthentication(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {

	md, exists := metadata.FromIncomingContext(ctx)

	if !exists {
		return nil, status.Errorf(codes.Unauthenticated, "metadata not provided")
	}

	values := md["authorization"]

	if len(values) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "authorization metadata not provided")
	}

	accessToken := values[0]

	conn, err := grpc.Dial(a.cfg.AuthServiceUrl())

	if err != nil {
		log.Fatalln("could not stablish connection to auth service, it may goes down: ", err.Error())
	}

	client := pb.NewAuthServiceClient(conn)

	in := &pb.AuthenticateRequest{Token: accessToken}

	if res, err := client.Authenticate(context.Background(), in); err != nil {
		return nil, err
	} else {
		log.Println("user ", res.Uid, " accessed resource ", info.FullMethod)
	}

	defer conn.Close()

	return handler(ctx, req)
}
