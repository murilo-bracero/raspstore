package handler

import (
	"net/http"
	"time"

	"github.com/murilo-bracero/raspstore/commons/pkg/logger"
	v1 "github.com/murilo-bracero/raspstore/idp/api/v1"
	u "github.com/murilo-bracero/raspstore/idp/internal/api/utils"
	"github.com/murilo-bracero/raspstore/idp/internal/usecase"
)

type LoginHandler interface {
	Login(w http.ResponseWriter, r *http.Request)
}

type loginHandler struct {
	loginUseCase usecase.LoginUseCase
}

func NewLoginHandler(loginUseCase usecase.LoginUseCase) LoginHandler {
	return &loginHandler{loginUseCase: loginUseCase}
}

func (c *loginHandler) Login(w http.ResponseWriter, r *http.Request) {
	logger.Info("Extracting credentials from header")

	username, password, ok := r.BasicAuth()

	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	logger.Info("Credentials extracted successfully")

	if err := r.ParseForm(); err != nil {
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
		return
	}

	responseType := r.PostForm.Get("response_type")
	mfaToken := r.PostForm.Get("mfa_token")

	tokenCredentials, err := c.loginUseCase.AuthenticateCredentials(username, password, mfaToken)

	if err != nil {
		logger.Error("Error while authenticating user: %s", err.Error())
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	if responseType == "code" {
		u.Send(w, v1.LoginResponse{
			AccessToken:  tokenCredentials.AccessToken,
			RefreshToken: tokenCredentials.RefreshToken,
		})
		return
	}

	if responseType == "token" {
		accessTokenCookie := createCookie("access_token", "Bearer "+tokenCredentials.AccessToken, tokenCredentials.ExpirestAt)
		http.SetCookie(w, accessTokenCookie)

		refreshTokenCookie := createCookie("refresh_token", tokenCredentials.RefreshToken, time.Time{})
		http.SetCookie(w, refreshTokenCookie)
	}
}

func createCookie(name string, value string, expiresAt time.Time) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    value,
		Expires:  expiresAt,
		Secure:   true,
		HttpOnly: true,
	}
}
