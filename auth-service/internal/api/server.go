package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/murilo-bracero/raspstore/auth-service/internal"
	"github.com/murilo-bracero/raspstore/auth-service/internal/api/handler"
	"github.com/murilo-bracero/raspstore/auth-service/internal/usecase"
)

func StartRestServer(ls usecase.LoginUseCase, puc usecase.GetUserUseCase) {
	loginHandler := handler.NewLoginHandler(ls)

	profileHandler := handler.NewProfileHandler(puc)

	router := NewRoutes(loginHandler, profileHandler).MountRoutes()
	http.Handle("/", router)
	log.Printf("Authentication API runing on port %d", internal.RestPort())
	http.ListenAndServe(fmt.Sprintf(":%d", internal.RestPort()), router)
}
