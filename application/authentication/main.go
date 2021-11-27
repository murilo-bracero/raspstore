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
	rp "raspstore.github.io/authentication/repository"
	ds "raspstore.github.io/authentication/repository/datastore"
	fb "raspstore.github.io/authentication/repository/firebase"
	mg "raspstore.github.io/authentication/repository/mongo"
	"raspstore.github.io/authentication/service"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Panicln("Could not load local variables")
	}

	cfg := db.NewConfig()

	var usersRepo rp.UsersRepository
	var credRepo rp.CredentialsRepository

	if cfg.UserDataStorage() == "mongodb" {
		conn, err := db.NewMongoConnection(context.Background(), cfg)
		usersRepo = mg.NewMongoUsersRepository(context.Background(), conn)

		if err != nil {
			log.Panicln(err)
		}

		defer conn.Close(context.Background())
	} else if cfg.UserDataStorage() == "datastore" {
		conn, err := db.NewDatastoreConnection(context.Background(), cfg)
		usersRepo = ds.NewDatastoreUsersRepository(context.Background(), conn)

		if err != nil {
			log.Panicln(err)
		}

		defer conn.Close()
	} else {
		log.Panicln("invalid user data storage option")
	}

	if cfg.CredentialsStorage() == "firebase" {
		conn, err := db.NewFirebaseConnection(context.Background())
		if err != nil {
			log.Panicln(err)
		}
		credRepo = fb.NewFireCredentials(context.Background(), conn)
	} else {
		log.Panicln("invalid credential storage option")
	}

	authService := service.NewAuthService(usersRepo, credRepo)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GrpcPort()))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterAuthServiceServer(grpcServer, authService)

	log.Printf("Authentication service running on [::]:%d\n", cfg.GrpcPort())

	grpcServer.Serve(lis)
}
