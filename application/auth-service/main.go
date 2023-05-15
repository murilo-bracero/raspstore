package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/joho/godotenv"
	"github.com/murilo-bracero/raspstore-protofiles/auth-service/pb"
	"google.golang.org/grpc"
	api "raspstore.github.io/auth-service/api"
	"raspstore.github.io/auth-service/api/controller"
	"raspstore.github.io/auth-service/db"
	gs "raspstore.github.io/auth-service/grpcservice"
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
	go startRestServer(&wg, loginService)
	wg.Wait()
}

func startRestServer(wg *sync.WaitGroup, ls usecase.LoginUseCase) {
	cc := controller.NewCredentialsController(ls)
	router := api.NewRoutes(cc).MountRoutes()
	http.Handle("/", router)
	log.Printf("Authentication API runing on port %d", db.RestPort())
	http.ListenAndServe(fmt.Sprintf(":%d", db.RestPort()), router)
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
