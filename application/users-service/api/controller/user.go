package controller

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/murilo-bracero/raspstore-protofiles/users-service/pb"
	"raspstore.github.io/users-service/api/dto"
	"raspstore.github.io/users-service/repository"
	"raspstore.github.io/users-service/utils"
)

type UserController interface {
	GetUser(w http.ResponseWriter, r *http.Request)
	ListUser(w http.ResponseWriter, r *http.Request)
}

type controller struct {
	repo         repository.UsersRepository
	usersService pb.UsersServiceServer
}

func NewUserController(repo repository.UsersRepository, us pb.UsersServiceServer) UserController {
	return &controller{repo: repo, usersService: us}
}

func (c *controller) GetUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id := params["id"]

	user, err := c.repo.FindById(id)

	if user == nil {
		w.WriteHeader(http.StatusNotFound)
		utils.Send(w, dto.ErrorResponse{Message: fmt.Sprintf("user with id %s does not exists", id), Code: "GU01"})
		return
	}

	if err != nil {
		w.WriteHeader(utils.ReqStatusCode(err))
		utils.Send(w, dto.ErrorResponse{Message: fmt.Sprintf("could not retrieve user with id %s", id), Reason: err.Error(), Code: "GU02"})
		return
	}

	utils.Send(w, user)
}

func (c *controller) ListUser(w http.ResponseWriter, r *http.Request) {

	users, err := c.repo.FindAll()

	if err != nil {
		w.WriteHeader(utils.ReqStatusCode(err))
		utils.Send(w, dto.ErrorResponse{Message: "could not retrieve users list", Reason: err.Error(), Code: "LU01"})
		return
	}

	utils.Send(w, users)
}
