package api

import (
	"fmt"
	"log"
	"net/http"

	"raspstore.github.io/auth-service/db"
	"raspstore.github.io/auth-service/internal/api/handler"
	"raspstore.github.io/auth-service/usecase"
)

func StartRestServer(ls usecase.LoginUseCase) {
	cc := handler.NewCredentialsHandler(ls)
	router := NewRoutes(cc).MountRoutes()
	http.Handle("/", router)
	log.Printf("Authentication API runing on port %d", db.RestPort())
	http.ListenAndServe(fmt.Sprintf(":%d", db.RestPort()), router)
}
