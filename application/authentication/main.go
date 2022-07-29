package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/joho/godotenv"
	"github.com/murilo-bracero/raspstore-protofiles/authentication/pb"
	"google.golang.org/grpc"
	api "raspstore.github.io/authentication/api"
	"raspstore.github.io/authentication/api/controller"
	"raspstore.github.io/authentication/db"
	interceptor "raspstore.github.io/authentication/interceptors"
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

	conn := initDatabase(cfg)

	defer conn.Close(context.Background())

	credRepo := initRepos(conn)

	tokenManager := token.NewTokenManager(cfg)

	authService := service.NewAuthService(credRepo, tokenManager)

	authInterceptor := interceptor.NewAuthInterceptor(tokenManager, credRepo)

	var wg sync.WaitGroup

	wg.Add(2)
	log.Println("bootstraping servers")
	go startGrpcServer(&wg, cfg, authInterceptor, authService)
	go startRestServer(&wg, cfg, authService)
	wg.Wait()
}

func startRestServer(wg *sync.WaitGroup, cfg db.Config, as pb.AuthServiceServer) {
	cc := controller.NewCredentialsController(as)
	router := api.NewRoutes(cc).MountRoutes()
	http.Handle("/", router)
	log.Printf("Authentication API runing on port %d", cfg.RestPort())
	http.ListenAndServe(fmt.Sprintf(":%d", cfg.RestPort()), router)
}

func startGrpcServer(wg *sync.WaitGroup, cfg db.Config, itc interceptor.AuthInterceptor, as pb.AuthServiceServer) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GrpcPort()))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(itc.WithUnaryAuthentication), grpc.StreamInterceptor(itc.WithStreamingAuthentication))
	pb.RegisterAuthServiceServer(grpcServer, as)

	log.Printf("Authentication service running on port:%d\n", cfg.GrpcPort())

	grpcServer.Serve(lis)
}

func initDatabase(cfg db.Config) db.MongoConnection {
	conn, err := db.NewMongoConnection(context.Background(), cfg)

	if err != nil {
		log.Panicln(err)
	}

	return conn
}

func initRepos(conn db.MongoConnection) rp.CredentialsRepository {
	return rp.NewCredentialsRepository(context.Background(), conn)
}
