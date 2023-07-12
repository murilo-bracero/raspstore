package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/murilo-bracero/raspstore/users-service/internal"
	"github.com/murilo-bracero/raspstore/users-service/internal/api/handler"
	"github.com/murilo-bracero/raspstore/users-service/internal/grpc"
	"github.com/murilo-bracero/raspstore/users-service/internal/service"
)

func StartRestServer(us service.UserService, ucs service.UserConfigService) {
	uh := handler.NewUserHandler(us, ucs)
	uch := handler.NewUserConfigHandler(ucs)
	auh := handler.NewAdminUserHandler(us, ucs)
	as := grpc.NewAuthService()
	router := NewRoutes(uh, uch, auh, as).MountRoutes()
	http.Handle("/", router)
	log.Printf("Authentication API runing on port %d", internal.RestPort())
	http.ListenAndServe(fmt.Sprintf(":%d", internal.RestPort()), router)
}
