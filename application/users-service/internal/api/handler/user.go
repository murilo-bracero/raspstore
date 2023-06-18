package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	v1 "raspstore.github.io/users-service/api/v1"
	"raspstore.github.io/users-service/internal"
	"raspstore.github.io/users-service/internal/api/middleware"
	u "raspstore.github.io/users-service/internal/api/utils"
	"raspstore.github.io/users-service/internal/model"
	"raspstore.github.io/users-service/internal/repository"
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
	userService          service.UserService
	userConfigRepository repository.UsersConfigRepository
}

func NewUserHandler(userService service.UserService, ucr repository.UsersConfigRepository) UserHandler {
	return &userHandler{userService: userService, userConfigRepository: ucr}
}

func (h *userHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	traceId := r.Context().Value(chiMiddleware.RequestIDKey).(string)
	cfg, err := h.userConfigRepository.Find()

	if err != nil {
		log.Printf("[ERROR] - [%s]: Could not retrieve configs: %s", traceId, err.Error())
		u.InternalServerError(w, traceId)
		return
	}

	if !cfg.AllowPublicUserCreation {
		handleBadRequest(w, "ERR006", "User creation not allowed", traceId)
		return
	}

	var req v1.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
		return
	}

	if err := validators.ValidateCreateUserRequest(req); err != nil {
		handleBadRequest(w, "ERR001", err.Error(), traceId)
		return
	}

	if len(req.Password) < cfg.MinPasswordLength {
		handleBadRequest(w, "ERR006", fmt.Sprintf("Password must be longer than %d", cfg.MinPasswordLength), traceId)
		return
	}

	usr := model.NewUserByCreateUserRequest(req)

	if err := h.userService.CreateUser(usr); err == internal.ErrUserAlreadyExists {
		log.Printf("[ERROR] - [%s]: User with [email=%s,username=%s] already exists in database", traceId, req.Email, req.Username)
		handleBadRequest(w, "ERR002", "User with provided email or username already exists", traceId)
		return
	} else if err != nil {
		log.Printf("[ERROR] - [%s]: Could not create user due to error: %s", traceId, err.Error())
		u.InternalServerError(w, traceId)
		return
	}

	log.Printf("[INFO] - [%s]: Created user succesfully. userId=%s", traceId, usr.UserId)

	u.Created(w, usr.ToUserResponse())
}

func (h *userHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	user, err := h.userService.GetUserById(id)

	if err == internal.ErrUserNotFound {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if err != nil {
		traceId := r.Context().Value(chiMiddleware.RequestIDKey).(string)
		log.Printf("[ERROR] - [%s]: Could search user with id %s in database: %s", traceId, id, err.Error())
		u.InternalServerError(w, traceId)
		return
	}

	u.Send(w, user)
}

func (h *userHandler) ListUser(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	size, _ := strconv.Atoi(r.URL.Query().Get("size"))

	if size == 0 || size > maxListSize {
		size = maxListSize
	}

	userPage, err := h.userService.GetAllUsersByPage(page, size)

	if err != nil {
		traceId := r.Context().Value(chiMiddleware.RequestIDKey).(string)
		log.Printf("[ERROR] - [%s]: Could list users due to error: %s", traceId, err.Error())
		u.InternalServerError(w, traceId)
		return
	}

	content := make([]v1.UserResponse, len(userPage.Content))
	for i, usr := range userPage.Content {
		content[i] = usr.ToUserResponse()
	}

	nextUrl := ""

	if len(content) == size {
		nextUrl = u.BuildPaginationNextUrl(r, page, size)
	}

	u.Send(w, v1.UserResponseList{
		Page:          page,
		Size:          size,
		TotalElements: userPage.Count,
		Next:          nextUrl,
		Content:       content,
	})
}

func (h *userHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	traceId := r.Context().Value(chiMiddleware.RequestIDKey).(string)
	id := chi.URLParam(r, "id")

	err := h.userService.RemoveUserById(id)

	if err != nil {
		log.Printf("[ERROR] - [%s]: Could not delete user with id=%s due to error: %s", traceId, id, err.Error())
		u.InternalServerError(w, traceId)
		return
	}

	requesterId := r.Context().Value(middleware.UserIdKey).(string)
	log.Printf("[ERROR] - [%s]: User removed successfully. requesterId=%s, userId=%s", traceId, requesterId, id)

	w.WriteHeader(http.StatusNoContent)
}

func (h *userHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	traceId := r.Context().Value(chiMiddleware.RequestIDKey).(string)

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
		http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
		return
	}

	if err != nil {
		log.Printf("[ERROR] - [%s]: Could not update user with id=%s due to error: %s", traceId, id, err.Error())
		u.InternalServerError(w, traceId)
		return
	}

	log.Printf("[ERROR] - [%s]: User updated successfully. userId=%s", traceId, id)

	u.Send(w, user.ToUserResponse())
}

func handleBadRequest(w http.ResponseWriter, code string, message string, traceId string) {
	u.BadRequest(w, v1.ErrorResponse{
		Code:    code,
		Message: message,
		TraceId: traceId,
	})
}
