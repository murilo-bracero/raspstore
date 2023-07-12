package grpc

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"raspstore.github.io/users-service/internal"
	"raspstore.github.io/users-service/internal/service"
	"raspstore.github.io/users-service/proto/v1/users-service/pb"
)

func StartGrpcServer(userConfigService service.UserConfigService) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", internal.GrpcPort()))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterUserConfigServiceServer(grpcServer, NewUserGrpcService(userConfigService))
	reflection.Register(grpcServer)

	log.Printf("UserConfig gRPC service running on port=%d", internal.GrpcPort())

	grpcServer.Serve(lis)
}
