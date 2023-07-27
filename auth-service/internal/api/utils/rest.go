package utils

import (
	"encoding/json"
	"io"
	"net/http"

	v1 "github.com/murilo-bracero/raspstore/auth-service/api/v1"
)

func HandleBadRequest(w http.ResponseWriter, code string, message string, traceId string) {
	BadRequest(w, v1.ErrorResponse{
		Code:    code,
		Message: message,
		TraceId: traceId,
	})
}

func BadRequest(w http.ResponseWriter, body v1.ErrorResponse) {
	w.WriteHeader(http.StatusBadRequest)
	Send(w, body)
}

func Send(w http.ResponseWriter, obj interface{}) {
	w.Header().Set("Content-Type", "application/json")
	jsonResponse, err := json.Marshal(obj)
	if err != nil {
		return
	}
	w.Write(jsonResponse)
}

func ParseBody(r io.Reader, v any) error {
	return json.NewDecoder(r).Decode(v)
}
