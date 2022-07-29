package controller

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/murilo-bracero/raspstore-protofiles/authentication/pb"
	"raspstore.github.io/authentication/api/dto"
)

type CredentialsController interface {
	Login(w http.ResponseWriter, r *http.Request)
}

type credsController struct {
	service pb.AuthServiceServer
}

func NewCredentialsController(service pb.AuthServiceServer) CredentialsController {
	return &credsController{service: service}
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
		send(w, nil)
		return
	}

	send(w, dto.LoginResponse{
		Token: res.Token,
	})

}
