package main

import (
	"context"
	"log"
	"sync"

	"github.com/joho/godotenv"
	"raspstore.github.io/users-service/internal/api"
	"raspstore.github.io/users-service/internal/database"
	"raspstore.github.io/users-service/internal/grpc"
	"raspstore.github.io/users-service/internal/repository"
	"raspstore.github.io/users-service/internal/service"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Could not load .env file. Using system variables instead")
	}

	ctx := context.Background()

	conn := initDatabase(ctx)

	defer conn.Close(ctx)

	usersRepository := repository.NewUsersRepository(ctx, conn)
	configRepository := repository.NewUsersConfigRepository(ctx, conn)

	userService := service.NewUserService(usersRepository, configRepository)
	userConfigService := service.NewUserConfigService(configRepository)

	var wg sync.WaitGroup

	wg.Add(2)
	log.Println("bootstraping servers")
	go grpc.StartGrpcServer(userConfigService)
	go api.StartRestServer(userService, userConfigService)
	wg.Wait()
}

func initDatabase(ctx context.Context) database.MongoConnection {
	conn, err := database.NewMongoConnection(ctx)

	if err != nil {
		log.Panicln(err)
	}

	return conn
}
