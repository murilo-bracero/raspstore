package api

import (
	"fmt"
	"log"
	"net/http"

	"raspstore.github.io/auth-service/internal"
	"raspstore.github.io/auth-service/internal/api/handler"
	"raspstore.github.io/auth-service/internal/service"
)

func StartRestServer(ls service.LoginService) {
	cc := handler.NewCredentialsHandler(ls)
	router := NewRoutes(cc).MountRoutes()
	http.Handle("/", router)
	log.Printf("Authentication API runing on port %d", internal.RestPort())
	http.ListenAndServe(fmt.Sprintf(":%d", internal.RestPort()), router)
}
