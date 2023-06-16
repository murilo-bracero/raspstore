package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	v1 "raspstore.github.io/users-service/api/v1"
	"raspstore.github.io/users-service/internal"
	"raspstore.github.io/users-service/internal/model"
	"raspstore.github.io/users-service/internal/service"
	"raspstore.github.io/users-service/internal/validators"
)

type AdminUserHandler interface {
	CreateUser(w http.ResponseWriter, r *http.Request)
}

type adminUserHandler struct {
	userService service.UserService
}

func NewAdminUserHandler(userService service.UserService) AdminUserHandler {
	return &adminUserHandler{userService: userService}
}

func (h *adminUserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req v1.AdminCreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
		return
	}

	if err := validators.ValidateCreateUserRequest(req.CreateUserRequest); err != nil {
		v1.BadRequest(w, v1.ErrorResponse{
			Code:    "ERR001",
			Message: err.Error(),
			TraceId: r.Context().Value(middleware.RequestIDKey).(string),
		})
		return
	}

	usr := model.NewUserByAdminCreateUserRequest(req)

	if err := h.userService.CreateUser(usr); err == internal.ErrUserAlreadyExists {
		traceId := r.Context().Value(middleware.RequestIDKey).(string)
		log.Printf("[ERROR] - [%s]: User with [email=%s,username=%s] already exists in database", traceId, req.Email, req.Username)
		v1.BadRequest(w, v1.ErrorResponse{
			Code:    "ERR002",
			Message: "User with provided email or username already exists",
			TraceId: traceId,
		})
		return
	} else if err != nil {
		traceId := r.Context().Value(middleware.RequestIDKey).(string)
		log.Printf("[ERROR] - [%s]: Could not create user due to error: %s", traceId, err.Error())
		v1.InternalServerError(w, traceId)
		return
	}

	v1.Created(w, usr.ToDto())
}