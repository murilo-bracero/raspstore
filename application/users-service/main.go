package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	api "raspstore.github.io/users-service/api"
	"raspstore.github.io/users-service/api/controller"
	"raspstore.github.io/users-service/api/middleware"
	"raspstore.github.io/users-service/db"
	rp "raspstore.github.io/users-service/repository"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Panicln("Could not load local variables")
	}

	cfg := db.NewConfig()

	conn := initDatabase(cfg)

	defer conn.Close(context.Background())

	usersRepo := initRepos(conn)

	authMiddleware := middleware.NewAuthMiddleware(cfg)

	startRestServer(cfg, usersRepo, authMiddleware)
}

func startRestServer(cfg db.Config, ur rp.UsersRepository, md middleware.AuthMiddleware) {
	uc := controller.NewUserController(ur)
	router := api.NewRoutes(uc).MountRoutes()
	log.Printf("Users Service API runing on port %d", cfg.RestPort())
	http.ListenAndServe(fmt.Sprintf(":%d", cfg.RestPort()), router)
}

func initDatabase(cfg db.Config) db.MongoConnection {
	conn, err := db.NewMongoConnection(context.Background(), cfg)

	if err != nil {
		log.Panicln(err)
	}

	return conn
}

func initRepos(conn db.MongoConnection) rp.UsersRepository {
	return rp.NewUsersRepository(context.Background(), conn)
}
