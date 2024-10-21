package handler

import (
	"net/http"
	"time"

	"github.com/murilo-bracero/raspstore/file-service/internal/domain/model"
)

func (l *Handler) Authenticate(w http.ResponseWriter, r *http.Request) {
	usr, pass, ok := r.BasicAuth()

	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	token, err := l.loginPAMUseCase.Execute(usr, pass)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	send(w, &model.LoginResponse{
		AccessToken: token,
		ExpiresIn:   int(time.Now().Add(1 * time.Hour).Unix()),
		Prefix:      "Bearer",
	})
}
