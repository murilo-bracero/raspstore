package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	u "raspstore.github.io/users-service/internal/api/utils"
	"raspstore.github.io/users-service/internal/service"
)

type UserConfigHandler interface {
	UpdateUserConfigs(w http.ResponseWriter, r *http.Request)
	GetUserConfigs(w http.ResponseWriter, r *http.Request)
}

type userConfigHandler struct {
	svc service.UserConfigService
}

func NewUserConfigHandler(svc service.UserConfigService) UserConfigHandler {
	return &userConfigHandler{svc: svc}
}

func (h *userConfigHandler) GetUserConfigs(w http.ResponseWriter, r *http.Request) {
	config, err := h.svc.GetUserConfig()

	if err != nil {
		traceId := r.Context().Value(middleware.RequestIDKey).(string)
		log.Printf("[ERROR] - [%s]: Could not retrieve user configuration: %s", traceId, err.Error())
		u.InternalServerError(w, traceId)
		return
	}

	u.Send(w, config.ToUserConfigurationResponse())
}

func (h *userConfigHandler) UpdateUserConfigs(w http.ResponseWriter, r *http.Request) {
	traceId := r.Context().Value(middleware.RequestIDKey).(string)
	config, err := h.svc.GetUserConfig()

	if err != nil {
		log.Printf("[ERROR] - [%s]: Could not retrieve user configuration: %s", traceId, err.Error())
		u.InternalServerError(w, traceId)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
		return
	}

	if err := h.svc.UpdateUserConfig(config); err != nil {
		log.Printf("[ERROR] - [%s]: Could not update user configuration: %s", traceId, err.Error())
		u.InternalServerError(w, traceId)
		return
	}

	log.Println("[INFO] - [%s]: Updated configs successfully", traceId)

	u.Send(w, config.ToUserConfigurationResponse())
}
