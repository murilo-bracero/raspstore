package client

import (
	"context"

	"github.com/murilo-bracero/raspstore/auth-service/proto/v1/auth-service/pb"
	"github.com/murilo-bracero/raspstore/commons/pkg/logger"
	"github.com/murilo-bracero/raspstore/commons/pkg/service"
	"github.com/murilo-bracero/raspstore/file-service/internal"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type authGrpcService struct{}

func NewAuthService() service.AuthService {
	return &authGrpcService{}
}

func (*authGrpcService) Authenticate(token string) (authResponse *pb.AuthenticateResponse, err error) {
	conn, err := grpc.Dial(internal.AuthServiceUrl(), grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		logger.Error("Could not stablish connection to auth service: %s", err.Error())
		return nil, err
	}

	defer conn.Close()

	client := pb.NewAuthServiceClient(conn)

	in := &pb.AuthenticateRequest{Token: token}

	return client.Authenticate(context.Background(), in)
}
