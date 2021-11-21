package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"raspstore.github.io/authentication/db"
	"raspstore.github.io/authentication/pb"
	"raspstore.github.io/authentication/repository"
	"raspstore.github.io/authentication/service"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Panicln("Could not load local variables")
	}

	cfg := db.NewConfig()

	var usersRepo repository.UsersRepository

	if cfg.UserDataStorage() == "mongodb" {
		conn, err := db.NewMongoConnection(context.Background(), cfg)
		usersRepo = repository.NewMongoUsersRepository(context.Background(), conn)

		if err != nil {
			log.Panicln(err)
		}

		defer conn.Close(context.Background())
	} else if cfg.UserDataStorage() == "datastore" {

	}

	authService := service.NewAuthService(usersRepo)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GrpcPort()))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterAuthServiceServer(grpcServer, authService)

	log.Printf("Authentication service running on [::]:%d\n", cfg.GrpcPort())

	grpcServer.Serve(lis)
}
