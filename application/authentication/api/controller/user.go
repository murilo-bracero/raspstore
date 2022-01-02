package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"raspstore.github.io/authentication/api/dto"
	"raspstore.github.io/authentication/model"
	"raspstore.github.io/authentication/pb"
	"raspstore.github.io/authentication/repository"
	"raspstore.github.io/authentication/service"
)

type UserController interface {
	SignUp(w http.ResponseWriter, r *http.Request)
	GetUser(w http.ResponseWriter, r *http.Request)
	DeleteUser(w http.ResponseWriter, r *http.Request)
	UpdateUser(w http.ResponseWriter, r *http.Request)
	ListUser(w http.ResponseWriter, r *http.Request)
}

type controller struct {
	repo    repository.UsersRepository
	service service.AuthService
}

func NewUserController(repo repository.UsersRepository, service service.AuthService) UserController {
	return &controller{repo: repo, service: service}
}

func (c *controller) SignUp(w http.ResponseWriter, r *http.Request) {

	var cUserReq dto.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&cUserReq); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	req := &pb.CreateUserRequest{
		Username:    cUserReq.Username,
		Password:    cUserReq.Password,
		Email:       cUserReq.Email,
		PhoneNumber: cUserReq.PhoneNumber,
	}

	if res, err := c.service.SignUp(context.Background(), req); err != nil {
		w.WriteHeader(reqStatusCode(err))
		send(w, dto.ErrorResponse{Message: "could not create user", Reason: err.Error(), Code: "RU01"})
	} else {
		w.WriteHeader(http.StatusCreated)
		send(w, model.User{
			UserId:      res.Id,
			Username:    res.Username,
			Email:       res.Email,
			PhoneNumber: res.PhoneNumber,
			CreatedAt:   time.Unix(res.CreatedAt, 0),
			UpdatedAt:   time.Unix(res.UpdatedAt, 0),
		})
	}

}

func (c *controller) GetUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id := params["id"]

	user, err := c.repo.FindById(id)

	if user == nil {
		w.WriteHeader(http.StatusNotFound)
		send(w, dto.ErrorResponse{Message: fmt.Sprintf("user with id %s does not exists", id), Code: "GU01"})
		return
	}

	if err != nil {
		w.WriteHeader(reqStatusCode(err))
		send(w, dto.ErrorResponse{Message: fmt.Sprintf("could not retrieve user with id %s", id), Reason: err.Error(), Code: "GU02"})
		return
	}

	send(w, user)
}

func (c *controller) DeleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id := params["id"]

	if err := c.repo.DeleteUser(id); err != nil {
		w.WriteHeader(reqStatusCode(err))
		send(w, dto.ErrorResponse{Message: fmt.Sprintf("could not delete user with id %s", id), Reason: err.Error(), Code: "DU01"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (c *controller) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var cUserReq dto.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&cUserReq); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	params := mux.Vars(r)

	id := params["id"]

	req := &pb.UpdateUserRequest{
		Id:          id,
		Username:    cUserReq.Username,
		Email:       cUserReq.Email,
		PhoneNumber: cUserReq.PhoneNumber,
	}

	if res, err := c.service.UpdateUser(context.Background(), req); err != nil {
		w.WriteHeader(reqStatusCode(err))
		send(w, dto.ErrorResponse{Message: fmt.Sprintf("could not update user with id %s", id), Reason: err.Error(), Code: "UU01"})
	} else {
		w.WriteHeader(reqStatusCode(err))
		send(w, model.User{
			UserId:      res.Id,
			Username:    res.Username,
			Email:       res.Email,
			PhoneNumber: res.PhoneNumber,
			CreatedAt:   time.Unix(res.CreatedAt, 0),
			UpdatedAt:   time.Unix(res.UpdatedAt, 0),
		})
	}
}

func (c *controller) ListUser(w http.ResponseWriter, r *http.Request) {

	users, err := c.repo.FindAll()

	if err != nil {
		w.WriteHeader(reqStatusCode(err))
		send(w, dto.ErrorResponse{Message: "could not retrieve users list", Reason: err.Error(), Code: "LU01"})
		return
	}

	send(w, users)
}
