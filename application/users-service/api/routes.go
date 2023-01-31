package api

import (
	"github.com/gorilla/mux"
	"raspstore.github.io/users-service/api/controller"
)

const userBaseRoute = "/users"

type Routes interface {
	MountRoutes() *mux.Router
}

type routes struct {
	userController controller.UserController
}

func NewRoutes(uc controller.UserController) Routes {
	return &routes{userController: uc}
}

func (r *routes) MountRoutes() *mux.Router {

	router := mux.NewRouter()

	router.HandleFunc(userBaseRoute+"/{id}", r.userController.GetUser).Methods("GET")
	router.HandleFunc(userBaseRoute, r.userController.ListUser).Methods("GET")

	return router

}
