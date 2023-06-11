package api

import (
	"fmt"
	"log"
	"net/http"

	"raspstore.github.io/users-service/internal"
	"raspstore.github.io/users-service/internal/api/handler"
	"raspstore.github.io/users-service/internal/repository"
)

func StartRestServer(ur repository.UsersRepository) {
	h := handler.NewUserHandler(ur)
	router := NewRoutes(h).MountRoutes()
	http.Handle("/", router)
	log.Printf("Authentication API runing on port %d", internal.RestPort())
	http.ListenAndServe(fmt.Sprintf(":%d", internal.RestPort()), router)
}
