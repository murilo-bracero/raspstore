package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"raspstore.github.io/users-service/api/controller"
)

const serviceBaseRoute = "/users-service"
const userBaseRoute = serviceBaseRoute + "/users"

type Routes interface {
	MountRoutes() *chi.Mux
}

type routes struct {
	userController controller.UserController
}

func NewRoutes(uc controller.UserController) Routes {
	return &routes{userController: uc}
}

func (rt *routes) MountRoutes() *chi.Mux {

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)

	router.Route(userBaseRoute, func(r chi.Router) {
		r.Get("/", rt.userController.ListUser)
		r.Post("/", rt.userController.CreateUser)
		r.Get("/{id}", rt.userController.GetUser)
		r.Put("/{id}", rt.userController.UpdateUser)
		r.Delete("/{id}", rt.userController.DeleteUser)
	})

	return router

}
