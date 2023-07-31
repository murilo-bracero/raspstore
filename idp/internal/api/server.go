package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/murilo-bracero/raspstore/idp/internal/api/handler"
	"github.com/murilo-bracero/raspstore/idp/internal/infra"
	"github.com/murilo-bracero/raspstore/idp/internal/usecase"
)

func StartRestServer(
	config *infra.Config,
	ls usecase.LoginUseCase,
	getUc usecase.GetUserUseCase,
	updateProfileUc usecase.UpdateProfileUseCase,
	updateUserUc usecase.UpdateUserUseCase,
	deleteUc usecase.DeleteUserUseCase,
	createUc usecase.CreateUserUseCase,
	listUs usecase.ListUsersUseCase) {

	loginHandler := handler.NewLoginHandler(ls)

	profileHandler := handler.NewProfileHandler(getUc, updateProfileUc, deleteUc)

	adminHandler := handler.NewAdminHandler(createUc, getUc, deleteUc, listUs, updateUserUc)

	router := NewRoutes(loginHandler, profileHandler, adminHandler, config).MountRoutes()
	http.Handle("/", router)
	log.Printf("Authentication API runing on port %d", config.RestPort)
	http.ListenAndServe(fmt.Sprintf(":%d", config.RestPort), router)
}
