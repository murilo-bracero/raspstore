package handler

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/murilo-bracero/raspstore/file-service/internal/domain/model"
)

func (l *Handler) Authenticate(w http.ResponseWriter, r *http.Request) {
	usr, pass, fine := r.BasicAuth()

	if !fine {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	token, err := l.LoginFunc(l.config, usr, pass)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	traceId := r.Context().Value(middleware.RequestIDKey).(string)

	ok(w, &model.LoginResponse{
		AccessToken: token,
		ExpiresIn:   int(time.Now().Add(1 * time.Hour).Unix()),
		Prefix:      "Bearer",
	}, traceId)
}
