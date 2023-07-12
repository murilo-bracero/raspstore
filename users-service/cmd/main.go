package main

import (
	"context"
	"log"
	"sync"

	"github.com/joho/godotenv"
	"github.com/murilo-bracero/raspstore/users-service/internal/api"
	"github.com/murilo-bracero/raspstore/users-service/internal/database"
	"github.com/murilo-bracero/raspstore/users-service/internal/grpc"
	"github.com/murilo-bracero/raspstore/users-service/internal/repository"
	"github.com/murilo-bracero/raspstore/users-service/internal/service"
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
