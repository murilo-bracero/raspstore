package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/murilo-bracero/raspstore/auth-service/internal"
	"github.com/murilo-bracero/raspstore/auth-service/internal/api/handler"
	"github.com/murilo-bracero/raspstore/auth-service/internal/service"
)

func StartRestServer(ls service.LoginService) {
	cc := handler.NewCredentialsHandler(ls)
	router := NewRoutes(cc).MountRoutes()
	http.Handle("/", router)
	log.Printf("Authentication API runing on port %d", internal.RestPort())
	http.ListenAndServe(fmt.Sprintf(":%d", internal.RestPort()), router)
}
