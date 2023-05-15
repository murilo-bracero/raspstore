package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	api "raspstore.github.io/users-service/api"
	"raspstore.github.io/users-service/api/controller"
	"raspstore.github.io/users-service/db"
	"raspstore.github.io/users-service/internal"
	rp "raspstore.github.io/users-service/repository"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Could not load .env file. Using system variables instead")
	}

	ctx := context.Background()

	conn := initDatabase(ctx)

	defer conn.Close(ctx)

	usersRepo := rp.NewUsersRepository(ctx, conn)

	startRestServer(usersRepo)
}

func startRestServer(ur rp.UsersRepository) {
	uc := controller.NewUserController(ur)
	router := api.NewRoutes(uc).MountRoutes()
	log.Printf("Users Service API runing on port %d", internal.RestPort())
	http.ListenAndServe(fmt.Sprintf(":%d", internal.RestPort()), router)
}

func initDatabase(ctx context.Context) db.MongoConnection {
	conn, err := db.NewMongoConnection(ctx)

	if err != nil {
		log.Panicln(err)
	}

	return conn
}
