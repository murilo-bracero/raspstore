package response

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/murilo-bracero/raspstore/file-service/internal/domain/model"
)

const traceIdHeaderKey = "X-Trace-Id"

func BadRequest(w http.ResponseWriter, body model.ErrorResponse, traceId string) {
	w.Header().Set(traceIdHeaderKey, traceId)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	send(w, body)
}

func InternalServerError(w http.ResponseWriter, traceId string) {
	w.Header().Set(traceIdHeaderKey, traceId)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func NotFound(w http.ResponseWriter, traceId string) {
	w.Header().Set(traceIdHeaderKey, traceId)
	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}

func UnprocessableEntity(w http.ResponseWriter, traceId string) {
	w.Header().Set(traceIdHeaderKey, traceId)
	http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
}

func Unauthorized(w http.ResponseWriter) {
	http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
}

func Created(w http.ResponseWriter, body interface{}, traceId string) {
	w.Header().Set(traceIdHeaderKey, traceId)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	send(w, body)
}

func Ok(w http.ResponseWriter, body interface{}, traceId string) {
	w.Header().Set(traceIdHeaderKey, traceId)
	send(w, body)
}

func send(w http.ResponseWriter, obj interface{}) {
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
