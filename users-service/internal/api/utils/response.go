package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	v1 "github.com/murilo-bracero/raspstore/users-service/api/v1"
)

func BadRequest(w http.ResponseWriter, body v1.ErrorResponse) {
	w.WriteHeader(http.StatusBadRequest)
	Send(w, body)
}

func InternalServerError(w http.ResponseWriter, traceId string) {
	w.WriteHeader(http.StatusInternalServerError)
	Send(w, v1.ErrorResponse{
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

func BuildPaginationNextUrl(r *http.Request, actualPage int, actualSize int) (nextUrl string) {
	resource := strings.Split(r.RequestURI, "?")[0]
	return fmt.Sprintf("%s%s?page=%d&size=%d", r.Host, resource, actualPage+1, actualSize)
}

func HandleBadRequest(w http.ResponseWriter, code string, message string, traceId string) {
	BadRequest(w, v1.ErrorResponse{
		Code:    code,
		Message: message,
		TraceId: traceId,
	})
}
