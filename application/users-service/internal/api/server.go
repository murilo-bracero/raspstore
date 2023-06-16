package api

import (
	"fmt"
	"log"
	"net/http"

	"raspstore.github.io/users-service/internal"
	"raspstore.github.io/users-service/internal/api/handler"
	"raspstore.github.io/users-service/internal/repository"
	"raspstore.github.io/users-service/internal/service"
)

func StartRestServer(us service.UserService, ucr repository.UsersConfigRepository) {
	uh := handler.NewUserHandler(us)
	uch := handler.NewUserConfigHandler(ucr)
	auh := handler.NewAdminUserHandler(us)
	router := NewRoutes(uh, uch, auh).MountRoutes()
	http.Handle("/", router)
	log.Printf("Authentication API runing on port %d", internal.RestPort())
	http.ListenAndServe(fmt.Sprintf(":%d", internal.RestPort()), router)
}
