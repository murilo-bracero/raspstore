package api

import (
	"github.com/gorilla/mux"
	"raspstore.github.io/file-manager/api/controller"
)

const fileBaseRoute = "/files"

type Routes interface {
	MountRoutes() *mux.Router
}

type routes struct {
	fc controller.FilesController
}

func NewRoutes(fc controller.FilesController) Routes {
	return &routes{fc: fc}
}

func (r *routes) MountRoutes() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc(fileBaseRoute, r.fc.Upload).Methods("POST")
	router.HandleFunc(fileBaseRoute+"/{id}", r.fc.Download).Methods("GET")
	router.HandleFunc(fileBaseRoute+"/{id}", r.fc.Update).Methods("PATCH")
	router.HandleFunc(fileBaseRoute+"/{id}", r.fc.Delete).Methods("DELETE")
	router.HandleFunc(fileBaseRoute, r.fc.ListFiles).Methods("GET")

	return router
}
