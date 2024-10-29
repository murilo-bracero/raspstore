package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/murilo-bracero/raspstore/file-service/internal/domain/model"
	"github.com/murilo-bracero/raspstore/file-service/internal/infra/validator"
)

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	traceId := r.Context().Value(middleware.RequestIDKey).(string)

	if !h.config.Auth.EnableUserRegister {
		slog.Warn("User registration is disabled", "traceId", traceId)
		unprocessableEntity(w, traceId)
		return
	}

	var req model.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		unprocessableEntity(w, traceId)
		return
	}

	if err := validator.ValidateRegisterRequest(&req); err != nil {
		badRequest(w, model.ErrorResponse{
			Message: err.Error(),
		}, traceId)
		return
	}

	if err := h.userFacade.Save(&req); err != nil {
		internalServerError(w, traceId)
		return
	}

	created(w, nil, traceId)
}
