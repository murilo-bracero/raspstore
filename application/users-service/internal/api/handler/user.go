package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	v1 "raspstore.github.io/users-service/api/v1"
	"raspstore.github.io/users-service/internal/model"
	"raspstore.github.io/users-service/internal/repository"
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

type handler struct {
	repo repository.UsersRepository
}

func NewUserHandler(repo repository.UsersRepository) UserHandler {
	return &handler{repo: repo}
}

func (c *handler) CreateUser(w http.ResponseWriter, r *http.Request) {
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

	if exists, err := c.repo.ExistsByEmailOrUsername(req.Email, req.Username); err != nil {
		traceId := r.Context().Value(middleware.RequestIDKey).(string)
		log.Printf("[ERROR] - [%s]: An unknown error occured while checking if user exists: %s", traceId, err.Error())
		v1.InternalServerError(w, traceId)
		return
	} else if exists {
		traceId := r.Context().Value(middleware.RequestIDKey).(string)
		log.Printf("[ERROR] - [%s]: User with [email=%s,username=%s] already exists in database", traceId, req.Email, req.Username)
		v1.BadRequest(w, v1.ErrorResponse{
			Code:    "ERR002",
			Message: "User with provided email or username already exists",
			TraceId: traceId,
		})
		return
	}

	usr := model.NewUserByCreateUserRequest(req)

	if hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost); err != nil {
		traceId := r.Context().Value(middleware.RequestIDKey).(string)
		log.Printf("[ERROR] - [%s]: Could not hash user password: %s", traceId, err.Error())
		v1.InternalServerError(w, traceId)
		return
	} else {
		usr.PasswordHash = string(hash)
	}

	if err := c.repo.Save(usr); err != nil {
		traceId := r.Context().Value(middleware.RequestIDKey).(string)
		log.Panicln(fmt.Sprintf("[ERROR] - [%s]: Could not create user due to error: %s", traceId, err.Error()))
		v1.InternalServerError(w, traceId)
		return
	}

	v1.Created(w, usr.ToDto())
}

func (c *handler) GetUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	user, err := c.repo.FindById(id)

	if err == mongo.ErrNoDocuments {
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

func (c *handler) ListUser(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	size, _ := strconv.Atoi(r.URL.Query().Get("size"))

	if size == 0 || size > maxListSize {
		size = maxListSize
	}

	userPage, err := c.repo.FindAll(page, size)

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

func (c *handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	err := c.repo.Delete(id)

	if err != nil {
		traceId := r.Context().Value(middleware.RequestIDKey).(string)
		v1.InternalServerError(w, traceId)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (c *handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var req v1.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
		return
	}

	id := chi.URLParam(r, "id")

	user, err := c.repo.FindById(id)

	if err == mongo.ErrNoDocuments {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if err != nil {
		traceId := r.Context().Value(middleware.RequestIDKey).(string)
		log.Printf("[ERROR] - [%s]: Could not search user with id %s in database: %s", traceId, id, err.Error())
		v1.InternalServerError(w, traceId)
		return
	}

	if req.Username != "" || req.Email != "" {
		if exists, err := c.repo.ExistsByEmailOrUsername(req.Email, req.Username); err != nil {
			traceId := r.Context().Value(middleware.RequestIDKey).(string)
			log.Printf("[ERROR] - [%s]: An unknown error occured while checking if user exists: %s", traceId, err.Error())
			v1.InternalServerError(w, traceId)
			return
		} else if exists {
			traceId := r.Context().Value(middleware.RequestIDKey).(string)
			log.Printf("[ERROR] - [%s]: User with [email=%s,username=%s] already exists in database", traceId, req.Email, req.Username)
			http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
			return
		}
	}

	if req.Email != "" {
		user.Email = req.Email
	}

	if req.Username != "" {
		user.Username = req.Username
	}

	if req.IsEnabled != nil {
		user.IsEnabled = *req.IsEnabled
	}

	if req.Password != "" {
		if hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost); err != nil {
			traceId := r.Context().Value(middleware.RequestIDKey).(string)
			log.Printf("[ERROR] - [%s]: Could not hash user password: %s", traceId, err.Error())
			v1.InternalServerError(w, traceId)
			return
		} else {
			user.PasswordHash = string(hash)
		}
	}

	if err := c.repo.Update(user); err != nil {
		traceId := r.Context().Value(middleware.RequestIDKey).(string)
		log.Printf("[ERROR] - [%s]: An unknown error occured while updating user with id %s: %s", traceId, id, err.Error())
		v1.InternalServerError(w, traceId)
		return
	}

	v1.Send(w, user.ToDto())
}
