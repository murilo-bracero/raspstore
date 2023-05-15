package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"raspstore.github.io/users-service/api/controller"
	mw "raspstore.github.io/users-service/api/middleware"
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
		r.With(mw.Authentication).Get("/", rt.userController.ListUser)
		r.With(mw.Authentication).Get("/{id}", rt.userController.GetUser)
		r.With(mw.Authentication).Put("/{id}", rt.userController.UpdateUser)
		r.With(mw.Authentication).Delete("/{id}", rt.userController.DeleteUser)
		r.Post("/", rt.userController.CreateUser)
	})

	return router

}
