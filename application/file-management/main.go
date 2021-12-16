package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"raspstore.github.io/file-manager/db"
	"raspstore.github.io/file-manager/pb"
	"raspstore.github.io/file-manager/repository"
	"raspstore.github.io/file-manager/service"
	"raspstore.github.io/file-manager/system"
)

func main() {
	ctx := context.Background()

	if err := godotenv.Load(); err != nil {
		log.Panicln("Could not load local variables")
	}

	cfg := db.NewConfig()

	conn, err := db.NewMongoConnection(context.Background(), cfg)

	if err != nil {
		log.Panicln(err)
	}

	defer conn.Close(context.Background())

	fileRepo := repository.NewFilesRepository(ctx, conn)
	diskStore := system.NewDiskStore(cfg.RootFolder())

	fileManagerService := service.NewFileManagerService(diskStore, fileRepo)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GrpcPort()))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterFileManagerServiceServer(grpcServer, fileManagerService)

	log.Printf("File Manager service running on [::]:%d\n", cfg.GrpcPort())

	grpcServer.Serve(lis)
}
