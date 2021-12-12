package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"raspstore.github.io/file-manager/db"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Panicln("Could not load local variables")
	}

	cfg := db.NewConfig()

	conn, err := db.NewMongoConnection(context.Background(), cfg)

	if err != nil {
		log.Panicln(err)
	}

	defer conn.Close(context.Background())

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GrpcPort()))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	log.Printf("Authentication service running on [::]:%d\n", cfg.GrpcPort())

	grpcServer.Serve(lis)
}
