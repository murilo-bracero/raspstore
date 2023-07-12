package grpc

import (
	"fmt"
	"log"
	"net"

	"github.com/murilo-bracero/raspstore-protofiles/auth-service/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"raspstore.github.io/auth-service/internal"
	"raspstore.github.io/auth-service/internal/service"
)

func StartGrpcServer(ts service.TokenService) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", internal.GrpcPort()))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterAuthServiceServer(grpcServer, NewAuthService(ts))
	reflection.Register(grpcServer)

	log.Printf("Authentication service running on port:%d\n", internal.GrpcPort())

	grpcServer.Serve(lis)
}
