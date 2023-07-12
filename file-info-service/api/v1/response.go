package v1

import (
	"encoding/json"
	"net/http"

	"github.com/murilo-bracero/raspstore/file-info-service/internal/model"
)

type FilePageResponse struct {
	Size          int                         `json:"size"`
	TotalElements int                         `json:"totalElements"`
	Page          int                         `json:"page"`
	Next          string                      `json:"next"`
	Content       []*model.FileMetadataLookup `json:"content"`
}

type ErrorResponse struct {
	Message string `json:"message,omitempty"`
	TraceId string `json:"traceId,omitempty"`
	Code    string `json:"code,omitempty"`
}

func BadRequest(w http.ResponseWriter, body ErrorResponse) {
	w.WriteHeader(http.StatusBadRequest)
	Send(w, body)
}

func InternalServerError(w http.ResponseWriter, traceId string) {
	w.WriteHeader(http.StatusInternalServerError)
	Send(w, ErrorResponse{
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
		w.Write(jsonResponse)
	}
}
