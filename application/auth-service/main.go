package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/joho/godotenv"
	"github.com/murilo-bracero/raspstore-protofiles/auth-service/pb"
	"google.golang.org/grpc"
	"raspstore.github.io/auth-service/db"
	gs "raspstore.github.io/auth-service/grpcservice"
	"raspstore.github.io/auth-service/internal/api"
	rp "raspstore.github.io/auth-service/repository"
	"raspstore.github.io/auth-service/token"
	"raspstore.github.io/auth-service/usecase"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Could not load .env file. Using system variables instead")
	}

	conn := initDatabase()

	defer conn.Close(context.Background())

	userRepository := initRepos(conn)

	tokenManager := token.NewTokenManager()

	authService := gs.NewAuthService(tokenManager)

	loginService := usecase.NewLoginUseCase(tokenManager, userRepository)

	var wg sync.WaitGroup

	wg.Add(2)
	log.Println("bootstraping servers")
	go startGrpcServer(&wg, authService)
	go api.StartRestServer(loginService)
	wg.Wait()
}

func startGrpcServer(wg *sync.WaitGroup, as pb.AuthServiceServer) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", db.GrpcPort()))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterAuthServiceServer(grpcServer, as)

	log.Printf("Authentication service running on port:%d\n", db.GrpcPort())

	grpcServer.Serve(lis)
}

func initDatabase() db.MongoConnection {
	conn, err := db.NewMongoConnection(context.Background())

	if err != nil {
		log.Panicln(err)
	}

	return conn
}

func initRepos(conn db.MongoConnection) rp.UsersRepository {
	return rp.NewUsersRepository(context.Background(), conn)
}
