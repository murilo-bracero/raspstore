package client

import (
	"context"
	"log"

	"github.com/murilo-bracero/raspstore/auth-service/proto/v1/auth-service/pb"
	"github.com/murilo-bracero/raspstore/commons/pkg/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"raspstore.github.io/fs-service/internal"
)

type authGrpcService struct {
	ctx context.Context
}

func NewAuthService(ctx context.Context) service.AuthService {
	return &authGrpcService{ctx: ctx}
}

func (a *authGrpcService) Authenticate(token string) (authResponse *pb.AuthenticateResponse, err error) {
	conn, err := grpc.Dial(internal.AuthServiceUrl(), grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Println("[ERROR] Could not stablish connection to auth service :", err.Error())
		return nil, err
	}

	defer conn.Close()

	client := pb.NewAuthServiceClient(conn)

	in := &pb.AuthenticateRequest{Token: token}

	return client.Authenticate(a.ctx, in)
}
