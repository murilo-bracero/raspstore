package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"raspstore.github.io/users-service/api/dto"
	"raspstore.github.io/users-service/model"
	"raspstore.github.io/users-service/repository"
	"raspstore.github.io/users-service/utils"
)

const defaultDateFormat = "2006-01-02 15:04:05"

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
		internalServerError(w, dto.ErrorResponse{
			Code:    "ERR099",
			Message: "Service currently unavailable",
			TraceId: traceId,
		})
		return
	}

	utils.Send(w, user)
}

func (c *controller) ListUser(w http.ResponseWriter, r *http.Request) {

	users, err := c.repo.FindAll()

	if err != nil {
		w.WriteHeader(utils.ReqStatusCode(err))
		utils.Send(w, dto.ErrorResponse{Message: "could not retrieve users list", Code: "LU01"})
		return
	}

	utils.Send(w, users)
}

func (c *controller) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	err := c.repo.Delete(id)

	if err != nil {
		w.WriteHeader(utils.ReqStatusCode(err))
		utils.Send(w, dto.ErrorResponse{Message: fmt.Sprintf("could not retrieve user with id %s", id), Code: "GU02"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
	utils.Send(w, nil)
}

func (c *controller) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
		return
	}

	if req.Email == "" || req.Username == "" {
		badRequest(w, dto.ErrorResponse{
			Code:    "ERR001",
			Message: "Field must not be null or empty",
			TraceId: r.Context().Value(middleware.RequestIDKey).(string),
		})
		return
	}

	if exists, err := c.repo.ExistsByEmailOrUsername(req.Email, req.Username); err != nil {
		traceId := r.Context().Value(middleware.RequestIDKey).(string)
		log.Printf("[ERROR] - [%s]: An unknown error occured while checking if user exists: %s", traceId, err.Error())
		internalServerError(w, dto.ErrorResponse{
			Code:    "ERR098",
			Message: "Service currently unavailable",
			TraceId: traceId,
		})
		return
	} else if exists {
		traceId := r.Context().Value(middleware.RequestIDKey).(string)
		log.Printf("[ERROR] - [%s]: User with [email=%s,username=%s] already exists in database", traceId, req.Email, req.Username)
		conflict(w, dto.ErrorResponse{
			Code:    "ERR002",
			Message: "User with provided email or username already exists",
			TraceId: traceId,
		})
		return
	}

	id := chi.URLParam(r, "id")

	usr := &model.User{
		UserId:      id,
		Username:    req.Username,
		Email:       req.Email,
		IsEnabled:   *req.IsEnabled,
		PhoneNumber: "",
	}

	if req.Password != "" {
		if hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost); err != nil {
			traceId := r.Context().Value(middleware.RequestIDKey).(string)
			log.Printf("[ERROR] - [%s]: Could not hash user password: %s", traceId, err.Error())
			internalServerError(w, dto.ErrorResponse{
				Code:    "ERR099",
				Message: "Service currently unavailable",
				TraceId: traceId,
			})
			return
		} else {
			usr.PasswordHash = string(hash)
		}
	}

	if err := c.repo.Update(usr); err != nil {
		traceId := r.Context().Value(middleware.RequestIDKey).(string)
		log.Printf("[ERROR] - [%s]: An unknown error occured while updating user with id %s: %s", traceId, id, err.Error())
		internalServerError(w, dto.ErrorResponse{
			Code:    "ERR098",
			Message: "Service currently unavailable",
			TraceId: traceId,
		})
		return
	}

	user, err := c.repo.FindById(id)

	if err != nil {
		traceId := r.Context().Value(middleware.RequestIDKey).(string)
		log.Printf("[ERROR] - [%s]: Could search user with id %s in database: %s", traceId, id, err.Error())
		internalServerError(w, dto.ErrorResponse{
			Code:    "ERR099",
			Message: "Service currently unavailable",
			TraceId: traceId,
		})
		return
	}

	utils.Send(w, user)
}

func (c *controller) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
		return
	}

	if req.Email == "" || req.Password == "" || req.Username == "" {
		badRequest(w, dto.ErrorResponse{
			Code:    "ERR001",
			Message: "Field must not be null or empty",
			TraceId: r.Context().Value(middleware.RequestIDKey).(string),
		})
		return
	}

	if exists, err := c.repo.ExistsByEmailOrUsername(req.Email, req.Username); err != nil {
		traceId := r.Context().Value(middleware.RequestIDKey).(string)
		log.Printf("[ERROR] - [%s]: An unknown error occured while checking if user exists: %s", traceId, err.Error())
		internalServerError(w, dto.ErrorResponse{
			Code:    "ERR098",
			Message: "Service currently unavailable",
			TraceId: traceId,
		})
		return
	} else if exists {
		traceId := r.Context().Value(middleware.RequestIDKey).(string)
		log.Printf("[ERROR] - [%s]: User with [email=%s,username=%s] already exists in database", traceId, req.Email, req.Username)
		badRequest(w, dto.ErrorResponse{
			Code:    "ERR002",
			Message: "User with provided email or username already exists",
			TraceId: traceId,
		})
		return
	}

	usr := &model.User{
		UserId:      uuid.NewString(),
		Username:    req.Username,
		Email:       req.Email,
		IsEnabled:   true,
		PhoneNumber: "",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost); err != nil {
		traceId := r.Context().Value(middleware.RequestIDKey).(string)
		log.Printf("[ERROR] - [%s]: Could not hash user password: %s", traceId, err.Error())
		internalServerError(w, dto.ErrorResponse{
			Code:    "ERR099",
			Message: "Service currently unavailable",
			TraceId: traceId,
		})
		return
	} else {
		usr.PasswordHash = string(hash)
	}

	if err := c.repo.Save(usr); err != nil {
		traceId := r.Context().Value(middleware.RequestIDKey).(string)
		log.Panicln(fmt.Sprintf("[ERROR] - [%s]: Could not create user due to error: %s", traceId, err.Error()))
		internalServerError(w, dto.ErrorResponse{
			Code:    "ERR002",
			Message: "Could not save user",
			TraceId: traceId,
		})
		return
	}

	created(w, dto.UserResponse{
		UserId:    usr.UserId,
		Username:  usr.Username,
		Email:     usr.Email,
		IsEnabled: usr.IsEnabled,
		CreatedAt: usr.CreatedAt.Format(defaultDateFormat),
		UpdatedAt: usr.UpdatedAt.Format(defaultDateFormat),
	})
}

func badRequest(w http.ResponseWriter, body dto.ErrorResponse) {
	w.WriteHeader(http.StatusBadRequest)
	utils.Send(w, body)
}

func conflict(w http.ResponseWriter, body dto.ErrorResponse) {
	w.WriteHeader(http.StatusConflict)
	utils.Send(w, body)
}

func internalServerError(w http.ResponseWriter, body dto.ErrorResponse) {
	w.WriteHeader(http.StatusInternalServerError)
	utils.Send(w, body)
}

func created(w http.ResponseWriter, body dto.UserResponse) {
	w.WriteHeader(http.StatusCreated)
	utils.Send(w, body)
}
