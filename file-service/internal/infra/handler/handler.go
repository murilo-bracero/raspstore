package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/murilo-bracero/raspstore/file-service/internal/application/facade"
	"github.com/murilo-bracero/raspstore/file-service/internal/application/usecase"
	"github.com/murilo-bracero/raspstore/file-service/internal/auth"
	"github.com/murilo-bracero/raspstore/file-service/internal/domain/model"
	"github.com/murilo-bracero/raspstore/file-service/internal/infra/config"
)

const traceIdHeaderKey = "X-Trace-Id"

type Handler struct {
	updateFileUseCase usecase.UpdateFileUseCase
	fileFacade        facade.FileFacade
	fileSystemFacade  facade.FileSystemFacade
	config            *config.Config
	// Public to make testing easier
	LoginFunc func(*config.Config, string, string) (string, error)
}

func New(
	updateFileUseCase usecase.UpdateFileUseCase,
	fileFacade facade.FileFacade,
	fileSystemFacade facade.FileSystemFacade,
	config *config.Config,

) *Handler {
	return &Handler{updateFileUseCase, fileFacade, fileSystemFacade, config, auth.LoginPAM}
}

func badRequest(w http.ResponseWriter, body model.ErrorResponse, traceId string) {
	w.Header().Set(traceIdHeaderKey, traceId)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	send(w, body)
}

func internalServerError(w http.ResponseWriter, traceId string) {
	w.Header().Set(traceIdHeaderKey, traceId)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func notFound(w http.ResponseWriter, traceId string) {
	w.Header().Set(traceIdHeaderKey, traceId)
	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}

func unprocessableEntity(w http.ResponseWriter, traceId string) {
	w.Header().Set(traceIdHeaderKey, traceId)
	http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
}

func created(w http.ResponseWriter, body interface{}, traceId string) {
	w.Header().Set(traceIdHeaderKey, traceId)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	send(w, body)
}

func ok(w http.ResponseWriter, body interface{}, traceId string) {
	w.Header().Set(traceIdHeaderKey, traceId)
	send(w, body)
}

func send(w http.ResponseWriter, obj interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if body, err := json.Marshal(obj); err == nil {
		if _, err := w.Write(body); err != nil {
			slog.Error("error while writing message to response body", "error", err)
		}
	}
}
