package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	v1 "github.com/murilo-bracero/raspstore/idp/api/v1"
	"github.com/murilo-bracero/raspstore/idp/internal"
	u "github.com/murilo-bracero/raspstore/idp/internal/api/utils"
	"github.com/murilo-bracero/raspstore/idp/internal/infra"
	"github.com/murilo-bracero/raspstore/idp/internal/model"
	"github.com/murilo-bracero/raspstore/idp/internal/usecase"
)

type AdminHandler interface {
	CreateUser(w http.ResponseWriter, r *http.Request)
	ListUsers(w http.ResponseWriter, r *http.Request)
	GetUserById(w http.ResponseWriter, r *http.Request)
	UpdateUserById(w http.ResponseWriter, r *http.Request)
	DeleteUser(w http.ResponseWriter, r *http.Request)
}

type adminHandler struct {
	config        *infra.Config
	createUseCase usecase.CreateUserUseCase
	getUseCase    usecase.GetUserUseCase
	deleteUseCase usecase.DeleteUserUseCase
	listUseCase   usecase.ListUsersUseCase
	updateUseCase usecase.UpdateUserUseCase
}

func NewAdminHandler(config *infra.Config,
	createUseCase usecase.CreateUserUseCase,
	getUseCase usecase.GetUserUseCase,
	deleteUseCase usecase.DeleteUserUseCase,
	listUseCase usecase.ListUsersUseCase,
	updateUseCase usecase.UpdateUserUseCase) AdminHandler {
	return &adminHandler{config: config,
		createUseCase: createUseCase,
		getUseCase:    getUseCase,
		deleteUseCase: deleteUseCase,
		listUseCase:   listUseCase,
		updateUseCase: updateUseCase}
}

func (h *adminHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	traceId := r.Context().Value(middleware.RequestIDKey).(string)
	var req v1.CreateUserRepresentation
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
		return
	}

	if err := validateCreateUserRequest(&req); err != nil {
		u.HandleBadRequest(w, "ERR001", err.Error(), traceId)
		return
	}

	if len(req.Password) < h.config.MinPasswordLength {
		u.HandleBadRequest(w, "ERR002", fmt.Sprintf("Password must have at least %d characters", h.config.MinPasswordLength), traceId)
		return
	}

	usr := model.NewUser(&req)

	if err := h.createUseCase.Execute(r.Context(), usr); err == internal.ErrConflict {
		u.HandleBadRequest(w, "ERR003", "User with provided email or username already exists", traceId)
		return
	} else if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	u.Created(w, usr.ToUserRepresentation())
}

func validateCreateUserRequest(req *v1.CreateUserRepresentation) error {
	if req.Username == "" {
		return internal.ErrInvalidUsername
	}

	if req.Password == "" {
		return internal.ErrInvalidPassword
	}

	if len(req.Roles) == 0 {
		return internal.ErrInvalidRoles
	}

	return nil
}

func (h *adminHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	size, _ := strconv.Atoi(r.URL.Query().Get("size"))
	username := r.URL.Query().Get("username")
	enabledQ := r.URL.Query().Get("enabled")

	var enabled *bool = nil
	if value, err := strconv.ParseBool(enabledQ); err == nil {
		enabled = &value
	}

	userPage, err := h.listUseCase.Execute(r.Context(), page, size, username, enabled)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	nextUrl := ""

	if len(userPage.Content) == size {
		nextUrl = u.BuildPaginationNextUrl(r, page, size)
	}

	u.Send(w, userPage.ToPageRepresentation(page, size, nextUrl))
}

func (h *adminHandler) GetUserById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "userId")

	user, err := h.getUseCase.Execute(r.Context(), id)

	if err == internal.ErrUserNotFound {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	u.Send(w, user.ToUserRepresentation())
}

func (h *adminHandler) UpdateUserById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "userId")

	var req v1.UpdateUserRepresentation
	if err := u.ParseBody(r.Body, &req); err != nil {
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
		return
	}

	if err := validateUpdateUserRepresentation(&req); err != nil {
		traceId := r.Context().Value(middleware.RequestIDKey).(string)
		u.HandleBadRequest(w, "ERR001", err.Error(), traceId)
		return
	}

	user := &model.User{
		UserId:       id,
		Username:     req.Username,
		IsEnabled:    req.Enabled,
		IsMfaEnabled: req.MfaEnabled,
		Permissions:  req.Roles,
	}

	if err := h.updateUseCase.Execute(r.Context(), user); err != nil {
		handleUpdateUserError(w, err)
		return
	}

	u.Send(w, user.ToUserRepresentation())
}

func validateUpdateUserRepresentation(req *v1.UpdateUserRepresentation) error {
	if req.Username == "" {
		return internal.ErrInvalidUsername
	}

	if len(req.Roles) == 0 {
		return internal.ErrInvalidRoles
	}

	return nil
}

func (h *adminHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "userId")

	if err := h.deleteUseCase.Execute(r.Context(), id); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	u.NoContent(w)
}
