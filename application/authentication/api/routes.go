package api

import (
	"github.com/gorilla/mux"
	"raspstore.github.io/authentication/api/controller"
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

	router.HandleFunc(userBaseRoute, r.userController.SignUp).Methods("POST")
	router.HandleFunc(userBaseRoute+"/{id}", r.userController.GetUser).Methods("GET")
	router.HandleFunc(userBaseRoute+"/{id}", r.userController.UpdateUser).Methods("PATCH")
	router.HandleFunc(userBaseRoute+"/{id}", r.userController.DeleteUser).Methods("DELETE")
	router.HandleFunc(userBaseRoute, r.userController.ListUser).Methods("GET")

	return router

}
