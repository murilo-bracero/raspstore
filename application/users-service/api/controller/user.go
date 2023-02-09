package controller

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
	"raspstore.github.io/users-service/api/dto"
	"raspstore.github.io/users-service/model"
	"raspstore.github.io/users-service/repository"
	"raspstore.github.io/users-service/validators"
)

const maxListSize = 50

type UserController interface {
	CreateUser(w http.ResponseWriter, r *http.Request)
	GetUser(w http.ResponseWriter, r *http.Request)
	ListUser(w http.ResponseWriter, r *http.Request)
	DeleteUser(w http.ResponseWriter, r *http.Request)
	UpdateUser(w http.ResponseWriter, r *http.Request)
}

type controller struct {
	repo repository.UsersRepository
}

func NewUserController(repo repository.UsersRepository) UserController {
	return &controller{repo: repo}
}

func (c *controller) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
		return
	}

	if err := validators.ValidateCreateUserRequest(req); err != nil {
		BadRequest(w, dto.ErrorResponse{
			Code:    "ERR001",
			Message: err.Error(),
			TraceId: r.Context().Value(middleware.RequestIDKey).(string),
		})
		return
	}

	if exists, err := c.repo.ExistsByEmailOrUsername(req.Email, req.Username); err != nil {
		traceId := r.Context().Value(middleware.RequestIDKey).(string)
		log.Printf("[ERROR] - [%s]: An unknown error occured while checking if user exists: %s", traceId, err.Error())
		InternalServerError(w, traceId)
		return
	} else if exists {
		traceId := r.Context().Value(middleware.RequestIDKey).(string)
		log.Printf("[ERROR] - [%s]: User with [email=%s,username=%s] already exists in database", traceId, req.Email, req.Username)
		BadRequest(w, dto.ErrorResponse{
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
		InternalServerError(w, traceId)
		return
	} else {
		usr.PasswordHash = string(hash)
	}

	if err := c.repo.Save(usr); err != nil {
		traceId := r.Context().Value(middleware.RequestIDKey).(string)
		log.Panicln(fmt.Sprintf("[ERROR] - [%s]: Could not create user due to error: %s", traceId, err.Error()))
		InternalServerError(w, traceId)
		return
	}

	Created(w, usr.ToDto())
}

func (c *controller) GetUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	user, err := c.repo.FindById(id)

	if err == mongo.ErrNoDocuments {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if err != nil {
		traceId := r.Context().Value(middleware.RequestIDKey).(string)
		log.Printf("[ERROR] - [%s]: Could search user with id %s in database: %s", traceId, id, err.Error())
		InternalServerError(w, traceId)
		return
	}

	Send(w, user)
}

func (c *controller) ListUser(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	size, _ := strconv.Atoi(r.URL.Query().Get("size"))

	if size == 0 || size > maxListSize {
		size = maxListSize
	}

	userPage, err := c.repo.FindAll(page, size)

	if err != nil {
		traceId := r.Context().Value(middleware.RequestIDKey).(string)
		log.Printf("[ERROR] - [%s]: Could list users due to error: %s", traceId, err.Error())
		InternalServerError(w, traceId)
		return
	}

	content := make([]dto.UserResponse, len(userPage.Content))
	for i, usr := range userPage.Content {
		content[i] = usr.ToDto()
	}

	nextUrl := ""

	if len(content) == size {
		nextUrl = fmt.Sprintf("%s/users-service/users?page=%d&size=%d", r.Host, page+1, size)
	}

	Send(w, dto.UserResponseList{
		Page:          page,
		Size:          size,
		TotalElements: userPage.Count,
		Next:          nextUrl,
		Content:       content,
	})
}

func (c *controller) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	err := c.repo.Delete(id)

	if err != nil {
		traceId := r.Context().Value(middleware.RequestIDKey).(string)
		InternalServerError(w, traceId)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (c *controller) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateUserRequest
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
		InternalServerError(w, traceId)
		return
	}

	if req.Username != "" || req.Email != "" {
		if exists, err := c.repo.ExistsByEmailOrUsername(req.Email, req.Username); err != nil {
			traceId := r.Context().Value(middleware.RequestIDKey).(string)
			log.Printf("[ERROR] - [%s]: An unknown error occured while checking if user exists: %s", traceId, err.Error())
			InternalServerError(w, traceId)
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
			InternalServerError(w, traceId)
			return
		} else {
			user.PasswordHash = string(hash)
		}
	}

	if err := c.repo.Update(user); err != nil {
		traceId := r.Context().Value(middleware.RequestIDKey).(string)
		log.Printf("[ERROR] - [%s]: An unknown error occured while updating user with id %s: %s", traceId, id, err.Error())
		InternalServerError(w, traceId)
		return
	}

	Send(w, user.ToDto())
}
