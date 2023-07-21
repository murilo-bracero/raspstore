package client

import (
	"context"

	"github.com/murilo-bracero/raspstore/commons/pkg/logger"
	"github.com/murilo-bracero/raspstore/file-info-service/internal"
	"github.com/murilo-bracero/raspstore/users-service/proto/v1/users-service/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type userConfigGrpcService struct{}

type UserConfigGrpcService interface {
	GetUserConfiguration() (*pb.UserConfiguration, error)
}

func NewUserConfigGrpcService() UserConfigGrpcService {
	return &userConfigGrpcService{}
}

func (u *userConfigGrpcService) GetUserConfiguration() (*pb.UserConfiguration, error) {
	conn, err := grpc.Dial(internal.UserServiceUrl(), grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		logger.Error("Could not stablish connection to user service: %s", err.Error())
		return nil, err
	}

	defer conn.Close()

	client := pb.NewUserConfigServiceClient(conn)

	in := &pb.GetUserConfigurationRequest{}

	return client.GetUserConfiguration(context.Background(), in)
}
