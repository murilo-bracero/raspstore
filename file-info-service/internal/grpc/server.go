package grpc

import (
	"fmt"
	"log"
	"net"

	"github.com/murilo-bracero/raspstore/file-info-service/internal"
	"github.com/murilo-bracero/raspstore/file-info-service/internal/grpc/server"
	"github.com/murilo-bracero/raspstore/file-info-service/internal/usecase"
	"github.com/murilo-bracero/raspstore/file-info-service/proto/v1/file-info-service/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func StartGrpcServer(gfuc usecase.GetFileUseCase, cfuc usecase.CreateFileUseCase) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", internal.GrpcPort()))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterFileInfoServiceServer(grpcServer, server.NewFileInfoService(gfuc, cfuc))
	reflection.Register(grpcServer)

	log.Printf("File Manager gRPC service running on [::]:%d\n", internal.GrpcPort())

	grpcServer.Serve(lis)
}
