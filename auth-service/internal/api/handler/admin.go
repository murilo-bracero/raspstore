package handler

import "net/http"

type AdminHandler interface {
	CreateUser(w http.ResponseWriter, r *http.Request)
	GetUserById(w http.ResponseWriter, r *http.Request)
	ListUsers(w http.ResponseWriter, r *http.Request)
	DeleteUser(w http.ResponseWriter, r *http.Request)
}

type adminHandler struct {
}

func NewAdminHandler() AdminHandler {
	return &adminHandler{}
}

func (h *adminHandler) CreateUser(w http.ResponseWriter, r *http.Request) {

}

func (h *adminHandler) GetUserById(w http.ResponseWriter, r *http.Request) {

}

func (h *adminHandler) ListUsers(w http.ResponseWriter, r *http.Request) {

}

func (h *adminHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {

}
