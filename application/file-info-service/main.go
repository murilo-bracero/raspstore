package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/joho/godotenv"
	"github.com/murilo-bracero/raspstore-protofiles/file-info-service/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"raspstore.github.io/file-manager/api"
	"raspstore.github.io/file-manager/api/controller"
	"raspstore.github.io/file-manager/db"
	"raspstore.github.io/file-manager/internal"
	"raspstore.github.io/file-manager/repository"
	"raspstore.github.io/file-manager/service"
)

func main() {
	ctx := context.Background()

	if err := godotenv.Load(); err != nil {
		log.Println("Could not load .env file. Using system variables instead")
	}

	conn, err := db.NewMongoConnection(context.Background())

	if err != nil {
		log.Panicln(err)
	}

	defer conn.Close(context.Background())

	fileRepo := repository.NewFilesRepository(ctx, conn)

	fileManagerService := service.NewFileInfoService(fileRepo)

	var wg sync.WaitGroup

	wg.Add(2)
	log.Println("bootstraping servers")
	go startGrpcServer(&wg, fileManagerService)
	go startRestServer(&wg, fileRepo)
	wg.Wait()
}

func startGrpcServer(wg *sync.WaitGroup, fileManagerService pb.FileInfoServiceServer) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", internal.GrpcPort()))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterFileInfoServiceServer(grpcServer, fileManagerService)
	reflection.Register(grpcServer)

	log.Printf("File Manager service running on [::]:%d\n", internal.GrpcPort())

	grpcServer.Serve(lis)
}

func startRestServer(wg *sync.WaitGroup, ur repository.FilesRepository) {
	fc := controller.NewFilesController(ur)
	router := api.NewRoutes(fc).MountRoutes()
	http.Handle("/", router)
	log.Printf("File Manager API runing on port %d", internal.RestPort())
	http.ListenAndServe(fmt.Sprintf(":%d", internal.RestPort()), router)
}
