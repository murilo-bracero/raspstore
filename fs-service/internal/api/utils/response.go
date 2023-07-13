package utils

import (
	"encoding/json"
	"net/http"

	v1 "raspstore.github.io/fs-service/api/v1"
)

func HandleBadRequest(w http.ResponseWriter, traceId string, message string, code string) {
	BadRequest(w, v1.ErrorResponse{
		TraceId: traceId,
		Message: message,
		Code:    code,
	})
}

func BadRequest(w http.ResponseWriter, body v1.ErrorResponse) {
	w.WriteHeader(http.StatusBadRequest)
	Send(w, body)
}

func Created(w http.ResponseWriter, body interface{}) {
	w.WriteHeader(http.StatusCreated)
	Send(w, body)
}

func Send(w http.ResponseWriter, obj interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if jsonResponse, err := json.Marshal(obj); err == nil {
		w.Write(jsonResponse)
	}
}
