package handler

import (
	"encoding/json"
	"log"
	"net/http"

	v1 "raspstore.github.io/auth-service/api/v1"
	"raspstore.github.io/auth-service/internal/service"
)

type CredentialsHandler interface {
	Login(w http.ResponseWriter, r *http.Request)
}

type credsHandler struct {
	loginService service.LoginService
}

func NewCredentialsHandler(ls service.LoginService) CredentialsHandler {
	return &credsHandler{loginService: ls}
}

func (c *credsHandler) Login(w http.ResponseWriter, r *http.Request) {
	log.Printf("initializing request")

	var lr v1.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&lr); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		v1.Send(w, nil)
		return
	}

	log.Printf("extracting authentication")

	username, password, ok := r.BasicAuth()

	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	if accessToken, refreshToken, err := c.loginService.AuthenticateCredentials(username, password, lr.MfaToken); err != nil {
		log.Printf("[ERROR] Error while authenticating user: %s", err.Error())
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	} else {
		v1.Send(w, v1.LoginResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		})
	}
}
