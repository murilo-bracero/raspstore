package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"raspstore.github.io/authentication/db"
	"raspstore.github.io/authentication/middleware"
	"raspstore.github.io/authentication/pb"
	rp "raspstore.github.io/authentication/repository"
	"raspstore.github.io/authentication/service"
	"raspstore.github.io/authentication/token"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Panicln("Could not load local variables")
	}

	cfg := db.NewConfig()

	conn, err := db.NewMongoConnection(context.Background(), cfg)

	if err != nil {
		log.Panicln(err)
	}

	defer conn.Close(context.Background())

	credRepo := rp.NewCredentialsRepository(context.Background(), conn)
	usersRepo := rp.NewMongoUsersRepository(context.Background(), conn)
	tokenManager := token.NewTokenManager(cfg)

	authService := service.NewAuthService(usersRepo, credRepo, tokenManager)

	authInterceptor := middleware.NewAuthInterceptor(usersRepo, tokenManager)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GrpcPort()))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(authInterceptor.WithAuthentication))
	pb.RegisterAuthServiceServer(grpcServer, authService)

	log.Printf("Authentication service running on [::]:%d\n", cfg.GrpcPort())

	grpcServer.Serve(lis)
}
