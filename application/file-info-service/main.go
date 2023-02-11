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
	"raspstore.github.io/file-manager/api"
	"raspstore.github.io/file-manager/api/controller"
	"raspstore.github.io/file-manager/api/middleware"
	"raspstore.github.io/file-manager/db"
	"raspstore.github.io/file-manager/repository"
	"raspstore.github.io/file-manager/service"
)

func main() {
	ctx := context.Background()

	if err := godotenv.Load(); err != nil {
		log.Println("Could not load .env file. Using system variables instead")
	}

	cfg := db.NewConfig()

	conn, err := db.NewMongoConnection(context.Background(), cfg)

	if err != nil {
		log.Panicln(err)
	}

	defer conn.Close(context.Background())

	fileRepo := repository.NewFilesRepository(ctx, conn)

	fileManagerService := service.NewFileManagerService(fileRepo)

	md := middleware.NewAuthMiddleware(cfg)

	var wg sync.WaitGroup

	wg.Add(2)
	log.Println("bootstraping servers")
	go startGrpcServer(&wg, cfg, fileManagerService)
	go startRestServer(&wg, cfg, fileRepo, md)
	wg.Wait()
}

func startGrpcServer(wg *sync.WaitGroup, cfg db.Config, fileManagerService pb.FileInfoServiceServer) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GrpcPort()))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterFileInfoServiceServer(grpcServer, fileManagerService)

	log.Printf("File Manager service running on [::]:%d\n", cfg.GrpcPort())

	grpcServer.Serve(lis)
}

func startRestServer(wg *sync.WaitGroup, cfg db.Config, ur repository.FilesRepository, md middleware.AuthMiddleware) {
	fc := controller.NewFilesController(ur)
	router := api.NewRoutes(fc).MountRoutes()
	http.Handle("/", router)
	log.Printf("File Manager API runing on port %d", cfg.RestPort())
	http.ListenAndServe(fmt.Sprintf(":%d", cfg.RestPort()), router)
}
