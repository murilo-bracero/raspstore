package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	api "raspstore.github.io/authentication/api"
	"raspstore.github.io/authentication/api/controller"
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
	usersRepo := rp.NewUsersRepository(context.Background(), conn)
	tokenManager := token.NewTokenManager(cfg)

	authService := service.NewAuthService(usersRepo, credRepo, tokenManager)

	authInterceptor := middleware.NewAuthInterceptor(usersRepo, tokenManager)

	go func() {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GrpcPort()))
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		grpcServer := grpc.NewServer(grpc.UnaryInterceptor(authInterceptor.WithAuthentication))
		pb.RegisterAuthServiceServer(grpcServer, authService)

		log.Printf("Authentication service running on [::]:%d\n", cfg.GrpcPort())

		grpcServer.Serve(lis)
	}()

	uc := controller.NewUserController(usersRepo)
	router := api.NewRoutes(uc).MountRoutes()
	http.Handle("/", router)
	log.Printf("Authentication API runing on port %d", cfg.RestPort())
	http.ListenAndServe(fmt.Sprintf(":%d", cfg.RestPort()), router)
}
