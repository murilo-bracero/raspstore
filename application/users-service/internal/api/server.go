package api

import (
	"fmt"
	"log"
	"net/http"

	"raspstore.github.io/users-service/internal"
	"raspstore.github.io/users-service/internal/api/handler"
	"raspstore.github.io/users-service/internal/grpc"
	"raspstore.github.io/users-service/internal/repository"
	"raspstore.github.io/users-service/internal/service"
)

func StartRestServer(us service.UserService, ucr repository.UsersConfigRepository) {
	uh := handler.NewUserHandler(us, ucr)
	uch := handler.NewUserConfigHandler(ucr)
	auh := handler.NewAdminUserHandler(us)
	as := grpc.NewAuthService()
	router := NewRoutes(uh, uch, auh, as).MountRoutes()
	http.Handle("/", router)
	log.Printf("Authentication API runing on port %d", internal.RestPort())
	http.ListenAndServe(fmt.Sprintf(":%d", internal.RestPort()), router)
}
