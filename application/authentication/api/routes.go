package api

import (
	"github.com/gorilla/mux"
	"raspstore.github.io/authentication/api/controller"
)

const userBaseRoute = "/users"

const loginRoute = "/auth"

type Routes interface {
	MountRoutes() *mux.Router
}

type routes struct {
	userController  controller.UserController
	credsController controller.CredentialsController
}

func NewRoutes(uc controller.UserController, cc controller.CredentialsController) Routes {
	return &routes{userController: uc, credsController: cc}
}

func (r *routes) MountRoutes() *mux.Router {

	router := mux.NewRouter()

	router.HandleFunc(userBaseRoute, r.userController.SignUp).Methods("POST")
	router.HandleFunc(userBaseRoute+"/{id}", r.userController.GetUser).Methods("GET")
	router.HandleFunc(userBaseRoute+"/{id}", r.userController.UpdateUser).Methods("PATCH")
	router.HandleFunc(userBaseRoute+"/{id}", r.userController.DeleteUser).Methods("DELETE")
	router.HandleFunc(userBaseRoute, r.userController.ListUser).Methods("GET")

	router.HandleFunc(loginRoute, r.credsController.Login).Methods("POST")

	return router

}
