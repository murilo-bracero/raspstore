package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	v1 "raspstore.github.io/users-service/api/v1"
	"raspstore.github.io/users-service/internal"
	"raspstore.github.io/users-service/internal/model"
	"raspstore.github.io/users-service/internal/service"
	"raspstore.github.io/users-service/internal/validators"
)

const maxListSize = 50

type UserHandler interface {
	CreateUser(w http.ResponseWriter, r *http.Request)
	GetUser(w http.ResponseWriter, r *http.Request)
	ListUser(w http.ResponseWriter, r *http.Request)
	DeleteUser(w http.ResponseWriter, r *http.Request)
	UpdateUser(w http.ResponseWriter, r *http.Request)
}

type userHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) UserHandler {
	return &userHandler{userService: userService}
}

func (h *userHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req v1.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
		return
	}

	if err := validators.ValidateCreateUserRequest(req); err != nil {
		v1.BadRequest(w, v1.ErrorResponse{
			Code:    "ERR001",
			Message: err.Error(),
			TraceId: r.Context().Value(middleware.RequestIDKey).(string),
		})
		return
	}

	usr := model.NewUserByCreateUserRequest(req)

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

func (h *userHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	user, err := h.userService.GetUserById(id)

	if err == internal.ErrUserNotFound {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if err != nil {
		traceId := r.Context().Value(middleware.RequestIDKey).(string)
		log.Printf("[ERROR] - [%s]: Could search user with id %s in database: %s", traceId, id, err.Error())
		v1.InternalServerError(w, traceId)
		return
	}

	v1.Send(w, user)
}

func (h *userHandler) ListUser(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	size, _ := strconv.Atoi(r.URL.Query().Get("size"))

	if size == 0 || size > maxListSize {
		size = maxListSize
	}

	userPage, err := h.userService.GetAllUsersByPage(page, size)

	if err != nil {
		traceId := r.Context().Value(middleware.RequestIDKey).(string)
		log.Printf("[ERROR] - [%s]: Could list users due to error: %s", traceId, err.Error())
		v1.InternalServerError(w, traceId)
		return
	}

	content := make([]v1.UserResponse, len(userPage.Content))
	for i, usr := range userPage.Content {
		content[i] = usr.ToDto()
	}

	nextUrl := ""

	if len(content) == size {
		nextUrl = fmt.Sprintf("%s/users-service/users?page=%d&size=%d", r.Host, page+1, size)
	}

	v1.Send(w, v1.UserResponseList{
		Page:          page,
		Size:          size,
		TotalElements: userPage.Count,
		Next:          nextUrl,
		Content:       content,
	})
}

func (h *userHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	err := h.userService.RemoveUserById(id)

	if err != nil {
		traceId := r.Context().Value(middleware.RequestIDKey).(string)
		log.Printf("[ERROR] - []: Could not delete user with id=%s due to error: %s", traceId, id, err.Error())
		v1.InternalServerError(w, traceId)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *userHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var req v1.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
		return
	}

	id := chi.URLParam(r, "id")

	user := &model.User{
		UserId:      id,
		Username:    req.Username,
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
	}

	user, err := h.userService.UpdateUser(user)

	if err == internal.ErrUserNotFound {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if err == internal.ErrEmailOrUsernameInUse {
		traceId := r.Context().Value(middleware.RequestIDKey).(string)
		log.Printf("[ERROR] - [%s]: User with [email=%s,username=%s] already exists in database", traceId, req.Email, req.Username)
		http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
		return
	}

	if err != nil {
		traceId := r.Context().Value(middleware.RequestIDKey).(string)
		log.Printf("[ERROR] - [%s]: Could not update user with id=%s due to error: %s", traceId, id, err.Error())
		v1.InternalServerError(w, traceId)
		return
	}

	v1.Send(w, user.ToDto())
}
