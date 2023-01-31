package controller

import (
	"encoding/json"
	"log"
	"net/http"

	"raspstore.github.io/auth-service/api/dto"
	"raspstore.github.io/auth-service/usecase"
	"raspstore.github.io/auth-service/utils"
)

type CredentialsController interface {
	Login(w http.ResponseWriter, r *http.Request)
}

type credsController struct {
	loginService usecase.LoginUseCase
}

func NewCredentialsController(ls usecase.LoginUseCase) CredentialsController {
	return &credsController{loginService: ls}
}

func (c *credsController) Login(w http.ResponseWriter, r *http.Request) {
	log.Printf("initializing request")

	var lr dto.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&lr); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		utils.Send(w, nil)
		return
	}

	log.Printf("extracting authentication")

	username, password, ok := r.BasicAuth()

	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	if accessToken, refreshToken, err := c.loginService.AuthenticateCredentials(username, password, lr.MfaToken); err != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	} else {
		utils.Send(w, dto.LoginResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		})
	}
}
