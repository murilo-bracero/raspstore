package handler

import (
	"net/http"
	"time"
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

	tokenCookie := http.Cookie{
		Name:     "JWT-TOKEN",
		Value:    token,
		Expires:  time.Now().Add(1 * time.Hour),
		Secure:   true,
		HttpOnly: true,
	}

	http.SetCookie(w, &tokenCookie)
}
