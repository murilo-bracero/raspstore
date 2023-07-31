package main

import (
	"context"
	"log"
	"sync"

	"github.com/joho/godotenv"
	"github.com/murilo-bracero/raspstore/idp/internal/api"
	"github.com/murilo-bracero/raspstore/idp/internal/database"
	"github.com/murilo-bracero/raspstore/idp/internal/grpc"
	"github.com/murilo-bracero/raspstore/idp/internal/infra"
	rp "github.com/murilo-bracero/raspstore/idp/internal/repository"
	"github.com/murilo-bracero/raspstore/idp/internal/usecase"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Could not load .env file. Using system variables instead")
	}

	config := infra.NewConfig()

	conn := initDatabase(config)

	defer conn.Close(context.Background())

	userRepository := initRepos(conn)

	loginUseCase := usecase.NewLoginUseCase(config, userRepository)
	getUserUseCase := usecase.NewGetUserUseCase(userRepository)
	updateProfileUseCase := usecase.NewUpdateProfileUseCase(userRepository)
	updateUserUseCase := usecase.NewUpdateUserUseCase(userRepository)
	deleteUseCase := usecase.NewDeleteUserUseCase(userRepository)
	createUseCase := usecase.NewCreateUserUseCase(userRepository, config)
	listUseCase := usecase.NewListUsersUseCase(userRepository)

	var wg sync.WaitGroup

	wg.Add(2)
	log.Println("bootstraping servers")
	go grpc.StartGrpcServer(config)
	go api.StartRestServer(config, loginUseCase, getUserUseCase, updateProfileUseCase, updateUserUseCase, deleteUseCase, createUseCase, listUseCase)
	wg.Wait()
}

func initDatabase(config *infra.Config) database.MongoConnection {
	conn, err := database.NewMongoConnection(context.Background(), config)

	if err != nil {
		log.Panicln(err)
	}

	return conn
}

func initRepos(conn database.MongoConnection) rp.UsersRepository {
	return rp.NewUsersRepository(context.Background(), conn)
}
