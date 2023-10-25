package utils

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/murilo-bracero/raspstore/file-service/internal/domain/model/response"
)

func HandleBadRequest(w http.ResponseWriter, traceId string, code string, message string) {
	BadRequest(w, response.ErrorResponse{
		TraceId: traceId,
		Code:    code,
		Message: message,
	})
}

func BadRequest(w http.ResponseWriter, body response.ErrorResponse) {
	w.WriteHeader(http.StatusBadRequest)
	Send(w, body)
}

func InternalServerError(w http.ResponseWriter, traceId string) {
	w.WriteHeader(http.StatusInternalServerError)
	Send(w, response.ErrorResponse{
		Code:    "ERR098",
		Message: "Service currently unavailable",
		TraceId: traceId,
	})
}

func Created(w http.ResponseWriter, body interface{}) {
	w.WriteHeader(http.StatusCreated)
	Send(w, body)
}

func Send(w http.ResponseWriter, obj interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if jsonResponse, err := json.Marshal(obj); err == nil {
		write(w, jsonResponse)
	}
}

func write(w http.ResponseWriter, body []byte) {
	if _, err := w.Write(body); err != nil {
		slog.Error("error while writing message to response body", "error", err)
	}
}
