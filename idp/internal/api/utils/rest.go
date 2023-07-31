package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	v1 "github.com/murilo-bracero/raspstore/idp/api/v1"
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

func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

func Created(w http.ResponseWriter, obj interface{}) {
	Send(w, obj)
	w.WriteHeader(http.StatusCreated)
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

func BuildPaginationNextUrl(r *http.Request, actualPage int, actualSize int) (nextUrl string) {
	resource := strings.Split(r.RequestURI, "?")[0]
	return fmt.Sprintf("%s%s?page=%d&size=%d", r.Host, resource, actualPage+1, actualSize)
}
