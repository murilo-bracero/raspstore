package api

import (
	"github.com/gorilla/mux"
	"raspstore.github.io/authentication/api/controller"
)

const loginRoute = "/auth"

type Routes interface {
	MountRoutes() *mux.Router
}

type routes struct {
	credsController controller.CredentialsController
}

func NewRoutes(cc controller.CredentialsController) Routes {
	return &routes{credsController: cc}
}

func (r *routes) MountRoutes() *mux.Router {

	router := mux.NewRouter()

	router.HandleFunc(loginRoute, r.credsController.Login).Methods("POST")

	return router

}
