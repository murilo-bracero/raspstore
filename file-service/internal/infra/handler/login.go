package handler

import "net/http"

type LoginHandler interface {
	Authenticate(w http.ResponseWriter, r *http.Request)
}

type loginHandler struct{}

func NewLoginHandler() *loginHandler {
	return &loginHandler{}
}
