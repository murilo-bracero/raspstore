package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	v1 "raspstore.github.io/users-service/api/v1"
	"raspstore.github.io/users-service/internal/model"
	"raspstore.github.io/users-service/internal/repository"
)

type UserConfigHandler interface {
	UpdateUserConfigs(w http.ResponseWriter, r *http.Request)
	GetUserConfigs(w http.ResponseWriter, r *http.Request)
}

type userConfigHandler struct {
	repository repository.UsersConfigRepository
}

func NewUserConfigHandler(repository repository.UsersConfigRepository) UserConfigHandler {
	return &userConfigHandler{repository: repository}
}

func (h *userConfigHandler) GetUserConfigs(w http.ResponseWriter, r *http.Request) {
	config, err := h.repository.Find()

	if err != nil {
		traceId := r.Context().Value(middleware.RequestIDKey).(string)
		log.Printf("[ERROR] - [%s]: Could not retrieve user configuration: %s", traceId, err.Error())
		v1.InternalServerError(w, traceId)
		return
	}

	v1.Send(w, toUserConfigurationResponse(config))
}

func (h *userConfigHandler) UpdateUserConfigs(w http.ResponseWriter, r *http.Request) {
	config, err := h.repository.Find()

	if err != nil {
		traceId := r.Context().Value(middleware.RequestIDKey).(string)
		log.Printf("[ERROR] - [%s]: Could not retrieve user configuration: %s", traceId, err.Error())
		v1.InternalServerError(w, traceId)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
		return
	}

	if err := h.repository.Update(config); err != nil {
		traceId := r.Context().Value(middleware.RequestIDKey).(string)
		log.Printf("[ERROR] - [%s]: Could not update user configuration: %s", traceId, err.Error())
		v1.InternalServerError(w, traceId)
		return
	}

	v1.Send(w, toUserConfigurationResponse(config))
}

func toUserConfigurationResponse(config *model.UserConfiguration) v1.UserConfigurationResponse {
	return v1.UserConfigurationResponse{
		StorageLimit:            config.StorageLimit,
		MinPasswordLength:       config.MinPasswordLength,
		AllowPublicUserCreation: config.AllowPublicUserCreation,
		AllowLoginWithEmail:     config.AllowLoginWithEmail,
		EnforceMfa:              config.EnforceMfa,
	}
}
