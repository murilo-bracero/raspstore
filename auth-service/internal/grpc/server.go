package grpc

import (
	"fmt"
	"log"
	"net"

	"github.com/murilo-bracero/raspstore/auth-service/internal/infra"
	"github.com/murilo-bracero/raspstore/auth-service/proto/v1/auth-service/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func StartGrpcServer(config *infra.Config) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", config.GrpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterAuthServiceServer(grpcServer, NewAuthService(config))
	reflection.Register(grpcServer)

	log.Printf("Authentication service running on port:%d\n", config.GrpcPort)

	grpcServer.Serve(lis)
}
