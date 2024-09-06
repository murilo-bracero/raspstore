package handler

import (
	"net/http"
	"time"

	"github.com/murilo-bracero/raspstore/file-service/internal/infra/response"
)

func (l *Handler) Authenticate(w http.ResponseWriter, r *http.Request) {
	usr, pass, ok := r.BasicAuth()

	if !ok {
		response.Unauthorized(w)
		return
	}

	token, err := l.loginPAMUseCase.Execute(usr, pass)

	if err != nil {
		response.Unauthorized(w)
		return
	}

	tokenCookie := http.Cookie{
		Name:     "JWT-TOKEN",
		Value:    token,
		Expires:  time.Now().Add(1 * time.Hour),
		Secure:   true,
		HttpOnly: true,
	}

	http.SetCookie(w, &tokenCookie)
}
