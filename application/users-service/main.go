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
	api "raspstore.github.io/users-service/api"
	"raspstore.github.io/users-service/api/controller"
	"raspstore.github.io/users-service/api/middleware"
	"raspstore.github.io/users-service/db"
	interceptor "raspstore.github.io/users-service/interceptors"
	"raspstore.github.io/users-service/pb"
	rp "raspstore.github.io/users-service/repository"
	"raspstore.github.io/users-service/service"
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

	usersService := service.NewUserService(usersRepo, credRepo)

	authInterceptor := interceptor.NewAuthInterceptor(cfg)

	authMiddleware := middleware.NewAuthMiddleware(cfg)

	var wg sync.WaitGroup

	wg.Add(2)
	log.Println("bootstraping servers")
	go startGrpcServer(&wg, cfg, authInterceptor, usersService)
	go startRestServer(&wg, cfg, usersRepo, usersService, authMiddleware)
	wg.Wait()
}

func startRestServer(wg *sync.WaitGroup, cfg db.Config, ur rp.UsersRepository, us pb.UsersServiceServer, md middleware.AuthMiddleware) {
	uc := controller.NewUserController(ur, us)
	router := api.NewRoutes(uc).MountRoutes()
	http.Handle("/", router)
	log.Printf("Users Service API runing on port %d", cfg.RestPort())
	http.ListenAndServe(fmt.Sprintf(":%d", cfg.RestPort()), md.Apply(router))
}

func startGrpcServer(wg *sync.WaitGroup, cfg db.Config, itc interceptor.AuthInterceptor, us pb.UsersServiceServer) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GrpcPort()))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(itc.WithUnaryAuthentication), grpc.StreamInterceptor(itc.WithStreamingAuthentication))
	pb.RegisterUsersServiceServer(grpcServer, us)

	log.Printf("Users Service gRPC running on [::]:%d\n", cfg.GrpcPort())

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
