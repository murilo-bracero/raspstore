package main

import (
	"context"
	"log"
	"sync"

	"github.com/joho/godotenv"
	"raspstore.github.io/file-manager/internal/api"
	db "raspstore.github.io/file-manager/internal/database"
	"raspstore.github.io/file-manager/internal/grpc"
	"raspstore.github.io/file-manager/internal/repository"
)

func main() {
	ctx := context.Background()

	if err := godotenv.Load(); err != nil {
		log.Println("Could not load .env file. Using system variables instead")
	}

	conn, err := db.NewMongoConnection(context.Background())

	if err != nil {
		log.Panicln(err)
	}

	defer conn.Close(context.Background())

	fileRepo := repository.NewFilesRepository(ctx, conn)

	var wg sync.WaitGroup

	wg.Add(2)
	log.Println("bootstraping servers")
	go grpc.StartGrpcServer(fileRepo)
	go api.StartApiServer(fileRepo)
	wg.Wait()
}
