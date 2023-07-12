package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	v1 "raspstore.github.io/users-service/api/v1"
	"raspstore.github.io/users-service/internal"
	u "raspstore.github.io/users-service/internal/api/utils"
	"raspstore.github.io/users-service/internal/model"
	"raspstore.github.io/users-service/internal/service"
	"raspstore.github.io/users-service/internal/validators"
)

type AdminUserHandler interface {
	CreateUser(w http.ResponseWriter, r *http.Request)
}

type adminUserHandler struct {
	userService       service.UserService
	userConfigService service.UserConfigService
}

func NewAdminUserHandler(userService service.UserService, userConfigService service.UserConfigService) AdminUserHandler {
	return &adminUserHandler{userService: userService, userConfigService: userConfigService}
}

func (h *adminUserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	traceId := r.Context().Value(middleware.RequestIDKey).(string)
	var req v1.AdminCreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
		return
	}

	if err := validators.ValidateCreateUserRequest(req.CreateUserRequest); err != nil {
		u.HandleBadRequest(w, "ERR001", err.Error(), traceId)
		return
	}

	usr := model.NewUserByAdminCreateUserRequest(req)

	if err := h.userConfigService.ValidateUser(usr, true); err != nil {
		u.HandleBadRequest(w, "ERR003", err.Error(), traceId)
		return
	}

	if err := h.userService.CreateUser(usr); err == internal.ErrUserAlreadyExists {
		traceId := r.Context().Value(middleware.RequestIDKey).(string)
		log.Printf("[ERROR] - [%s]: User with [email=%s,username=%s] already exists in database", traceId, req.Email, req.Username)
		u.HandleBadRequest(w, "ERR002", "User with provided email or username already exists", traceId)
		return
	} else if err != nil {
		log.Printf("[ERROR] - [%s]: Could not create user due to error: %s", traceId, err.Error())
		u.InternalServerError(w, traceId)
		return
	}

	u.Created(w, usr.ToAdminUserResponse())
}
