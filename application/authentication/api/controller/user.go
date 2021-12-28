package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	api "raspstore.github.io/authentication/api/dto"
	"raspstore.github.io/authentication/model"
	"raspstore.github.io/authentication/repository"
)

type UserController interface {
	SignUp(w http.ResponseWriter, r *http.Request)
	GetUser(w http.ResponseWriter, r *http.Request)
	DeleteUser(w http.ResponseWriter, r *http.Request)
	UpdateUser(w http.ResponseWriter, r *http.Request)
	ListUser(w http.ResponseWriter, r *http.Request)
}

type controller struct {
	userRepository repository.UsersRepository
}

func NewUserController(userRepository repository.UsersRepository) UserController {
	return &controller{userRepository: userRepository}
}

func (c *controller) SignUp(w http.ResponseWriter, r *http.Request) {
	var cUserReq api.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&cUserReq); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user := new(model.User)
	user.FromCreateRequest(cUserReq)

	c.userRepository.Save(user)

	w.WriteHeader(http.StatusOK)
	send(w, user)
}

func (c *controller) GetUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id := params["id"]

	user, err := c.userRepository.FindById(id)

	if user == nil {
		w.WriteHeader(http.StatusNotFound)
		er := new(api.ErrorResponse)
		er.Message = fmt.Sprintf("user with id %s does not exists", id)
		send(w, er)
		return
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		er := new(api.ErrorResponse)
		er.Message = "could not retrieve user with id " + id
		er.Reason = err.Error()
		er.Code = "GU01"
		send(w, er)
		return
	}

	send(w, user)
}

func (c *controller) DeleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id := params["id"]

	if err := c.userRepository.DeleteUser(id); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		er := new(api.ErrorResponse)
		er.Message = "could not delete user with id " + id
		er.Reason = err.Error()
		er.Code = "DU01"
		send(w, er)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (c *controller) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var cUserReq api.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&cUserReq); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	params := mux.Vars(r)

	id := params["id"]

	user := new(model.User)
	user.FromUpdateRequest(cUserReq, id)

	if err := c.userRepository.UpdateUser(user); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		er := new(api.ErrorResponse)
		er.Message = "could not update user with id " + id
		er.Reason = err.Error()
		er.Code = "UU01"
		send(w, er)
		return
	}

	send(w, user)

}

func (c *controller) ListUser(w http.ResponseWriter, r *http.Request) {

	users, err := c.userRepository.FindAll()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		er := new(api.ErrorResponse)
		er.Message = "could not retrieve users list"
		er.Reason = err.Error()
		er.Code = "LU01"
		send(w, er)
		return
	}

	send(w, users)
}

func send(w http.ResponseWriter, obj interface{}) {
	w.Header().Set("Content-Type", "application/json")
	jsonResponse, err := json.Marshal(obj)
	if err != nil {
		return
	}
	w.Write(jsonResponse)
}
