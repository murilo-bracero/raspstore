package grpc

import (
	"fmt"
	"log"
	"net"

	"github.com/murilo-bracero/raspstore/users-service/internal"
	"github.com/murilo-bracero/raspstore/users-service/internal/service"
	"github.com/murilo-bracero/raspstore/users-service/proto/v1/users-service/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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
