package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	api "raspstore.github.io/authentication/api"
	"raspstore.github.io/authentication/api/controller"
	"raspstore.github.io/authentication/api/middleware"
	"raspstore.github.io/authentication/db"
	interceptor "raspstore.github.io/authentication/interceptors"
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

	conn := initDatabase(cfg)

	defer conn.Close(context.Background())

	credRepo, usersRepo := initRepos(conn)

	tokenManager := token.NewTokenManager(cfg)

	authService := service.NewAuthService(usersRepo, credRepo, tokenManager)

	authInterceptor := interceptor.NewAuthInterceptor(usersRepo, tokenManager)

	authMiddleware := middleware.NewAuthMiddleware(tokenManager)

	var wg sync.WaitGroup

	wg.Add(2)
	log.Println("bootstraping servers")
	go startGrpcServer(&wg, cfg, authInterceptor, authService)
	go startRestServer(&wg, cfg, usersRepo, authService, authMiddleware)
	wg.Wait()
}

func startRestServer(wg *sync.WaitGroup, cfg db.Config, ur rp.UsersRepository, svc service.AuthService, md middleware.AuthMiddleware) {
	uc := controller.NewUserController(ur, svc)
	cc := controller.NewCredentialsController(ur, svc)
	router := api.NewRoutes(uc, cc).MountRoutes()
	http.Handle("/", router)
	log.Printf("Authentication API runing on port %d", cfg.RestPort())
	http.ListenAndServe(fmt.Sprintf(":%d", cfg.RestPort()), md.Apply(router))
}

func startGrpcServer(wg *sync.WaitGroup, cfg db.Config, itc interceptor.AuthInterceptor, svc pb.AuthServiceServer) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GrpcPort()))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(itc.WithAuthentication))
	pb.RegisterAuthServiceServer(grpcServer, svc)

	log.Printf("Authentication service running on [::]:%d\n", cfg.GrpcPort())

	grpcServer.Serve(lis)
}

func initDatabase(cfg db.Config) db.MongoConnection {
	conn, err := db.NewMongoConnection(context.Background(), cfg)

	if err != nil {
		log.Panicln(err)
	}

	return conn
}

func initRepos(conn db.MongoConnection) (rp.CredentialsRepository, rp.UsersRepository) {
	return rp.NewCredentialsRepository(context.Background(), conn),
		rp.NewUsersRepository(context.Background(), conn)
}
