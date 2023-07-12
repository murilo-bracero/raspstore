package main

import (
	"context"
	"log"
	"sync"

	"github.com/joho/godotenv"
	"github.com/murilo-bracero/raspstore/auth-service/internal/api"
	"github.com/murilo-bracero/raspstore/auth-service/internal/database"
	"github.com/murilo-bracero/raspstore/auth-service/internal/grpc"
	rp "github.com/murilo-bracero/raspstore/auth-service/internal/repository"
	"github.com/murilo-bracero/raspstore/auth-service/internal/service"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Could not load .env file. Using system variables instead")
	}

	conn := initDatabase()

	defer conn.Close(context.Background())

	userRepository := initRepos(conn)

	tokenService := service.NewTokenService()

	loginService := service.NewLoginService(tokenService, userRepository)

	var wg sync.WaitGroup

	wg.Add(2)
	log.Println("bootstraping servers")
	go grpc.StartGrpcServer(tokenService)
	go api.StartRestServer(loginService)
	wg.Wait()
}

func initDatabase() database.MongoConnection {
	conn, err := database.NewMongoConnection(context.Background())

	if err != nil {
		log.Panicln(err)
	}

	return conn
}

func initRepos(conn database.MongoConnection) rp.UsersRepository {
	return rp.NewUsersRepository(context.Background(), conn)
}
