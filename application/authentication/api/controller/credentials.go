package controller

import (
	"context"
	"encoding/json"
	"net/http"

	"raspstore.github.io/authentication/api/dto"
	"raspstore.github.io/authentication/pb"
	"raspstore.github.io/authentication/repository"
)

type CredentialsController interface {
	Login(w http.ResponseWriter, r *http.Request)
}

type credsController struct {
	repo    repository.UsersRepository
	service pb.AuthServiceServer
}

func NewCredentialsController(repo repository.UsersRepository, service pb.AuthServiceServer) CredentialsController {
	return &credsController{repo: repo, service: service}
}

func (c *credsController) Login(w http.ResponseWriter, r *http.Request) {
	var lr dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&lr); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	req := &pb.LoginRequest{
		Email:    lr.Email,
		Password: lr.Password,
	}

	res, err := c.service.Login(context.Background(), req)

	if err != nil {
		w.WriteHeader(reqStatusCode(err))
		send(w, dto.ErrorResponse{Message: "could not make login with given credentials", Reason: err.Error(), Code: "LG01"})
		return
	}

	usr, err := c.repo.FindByEmail(req.Email)
	if err != nil {
		w.WriteHeader(reqStatusCode(err))
		send(w, dto.ErrorResponse{Message: "could not retrieve logged in user", Reason: err.Error(), Code: "LG02"})
		return
	}

	send(w, dto.LoginResponse{
		Token: res.Token,
		User:  *usr,
	})

}
