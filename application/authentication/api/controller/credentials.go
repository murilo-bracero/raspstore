package controller

import (
	"context"
	"encoding/json"
	"net/http"

	api "raspstore.github.io/authentication/api/dto"
	"raspstore.github.io/authentication/pb"
	"raspstore.github.io/authentication/repository"
	"raspstore.github.io/authentication/service"
)

type CredentialsController interface {
	Login(w http.ResponseWriter, r *http.Request)
}

type credsController struct {
	repo    repository.UsersRepository
	service service.AuthService
}

func NewCredentialsController(repo repository.UsersRepository, service service.AuthService) CredentialsController {
	return &credsController{repo: repo, service: service}
}

func (c *credsController) Login(w http.ResponseWriter, r *http.Request) {
	var lr api.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&lr); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	req := &pb.LoginRequest{
		Email:    lr.Email,
		Password: lr.Password,
	}

	res, err := c.service.Login(context.Background(), req)
	w.WriteHeader(reqStatusCode(err))

	if err != nil {
		er := new(api.ErrorResponse)
		er.Message = "could not make login with given credentials"
		er.Reason = err.Error()
		er.Code = "LG01"
		send(w, er)
		return
	}

	usr, err := c.repo.FindByEmail(req.Email)
	w.WriteHeader(reqStatusCode(err))
	if err != nil {
		er := new(api.ErrorResponse)
		er.Message = "could not retrieve logged in user"
		er.Reason = err.Error()
		er.Code = "LG02"
		send(w, er)
		return
	}

	send(w, api.LoginResponse{
		Token: res.Token,
		User:  *usr,
	})

}
