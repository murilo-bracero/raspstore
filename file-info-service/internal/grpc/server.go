package grpc

import (
	"fmt"
	"log"
	"net"

	"github.com/murilo-bracero/raspstore-protofiles/file-info-service/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"raspstore.github.io/file-manager/internal"
	"raspstore.github.io/file-manager/internal/repository"
)

func StartGrpcServer(fileRepository repository.FilesRepository) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", internal.GrpcPort()))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterFileInfoServiceServer(grpcServer, NewFileInfoService(fileRepository))
	reflection.Register(grpcServer)

	log.Printf("File Manager gRPC service running on [::]:%d\n", internal.GrpcPort())

	grpcServer.Serve(lis)
}
