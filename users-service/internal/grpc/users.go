package grpc

import (
	"context"
	"log"

	"raspstore.github.io/users-service/internal/service"
	"raspstore.github.io/users-service/proto/v1/users-service/pb"
)

type userConfigGrpcService struct {
	userConfigService service.UserConfigService
	pb.UnimplementedUserConfigServiceServer
}

func NewUserGrpcService(ucs service.UserConfigService) pb.UserConfigServiceServer {
	return &userConfigGrpcService{userConfigService: ucs}
}

func (s *userConfigGrpcService) GetUserConfiguration(ctx context.Context, req *pb.GetUserConfigurationRequest) (response *pb.UserConfiguration, error_ error) {
	config, error_ := s.userConfigService.GetUserConfig()

	if error_ != nil {
		log.Printf("[ERROR] - Could not get user configurations: %s", error_.Error())
		return
	}

	return &pb.UserConfiguration{
		StorageLimit:            config.StorageLimit,
		MinPasswordLength:       int64(config.MinPasswordLength),
		AllowPublicUserCreation: config.AllowPublicUserCreation,
		AllowLoginWithEmail:     config.AllowLoginWithEmail,
		EnforceMfa:              config.EnforceMfa,
	}, nil
}
