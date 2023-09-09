package handler

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/lestrrat-go/jwx/jwt"
	v1 "github.com/murilo-bracero/raspstore/file-service/api/v1"
	m "github.com/murilo-bracero/raspstore/file-service/internal/api/middleware"
	u "github.com/murilo-bracero/raspstore/file-service/internal/api/utils"
	"github.com/murilo-bracero/raspstore/file-service/internal/converter"
	"github.com/murilo-bracero/raspstore/file-service/internal/infra"
	"github.com/murilo-bracero/raspstore/file-service/internal/model"
	"github.com/murilo-bracero/raspstore/file-service/internal/usecase"
)

type UploadHandler interface {
	Upload(w http.ResponseWriter, r *http.Request)
}

type uploadHandler struct {
	config            *infra.Config
	uploadUseCase     usecase.UploadFileUseCase
	createFileUseCase usecase.CreateFileUseCase
}

func NewUploadHandler(config *infra.Config, uploadUseCase usecase.UploadFileUseCase, createFileUseCase usecase.CreateFileUseCase) UploadHandler {
	return &uploadHandler{config: config, uploadUseCase: uploadUseCase, createFileUseCase: createFileUseCase}
}

func (h *uploadHandler) Upload(w http.ResponseWriter, r *http.Request) {
	usr := r.Context().Value(m.UserClaimsCtxKey).(jwt.Token)
	traceId := r.Context().Value(middleware.RequestIDKey).(string)
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		slog.Error("Could not allocate MultipartForm parser", "traceId", traceId, "error", err)
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
		return
	}

	file, header, err := r.FormFile("file")

	if err != nil {
		slog.Error("Could not open Multipart Form file", "traceId", traceId, "error", err)
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
		return
	}

	defer file.Close()

	fm := converter.CreateFile(header.Filename, header.Size, false, usr.Subject())

	if err := h.uploadUseCase.Execute(r.Context(), fm, file); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if err := h.createFileUseCase.Execute(fm); err != nil {
		h.handleCreateUseCaseError(w, fm)
		return
	}

	u.Created(w, &v1.UploadSuccessResponse{
		FileId:   fm.FileId,
		Filename: fm.Filename,
		OwnerId:  fm.Owner,
	})
}

func (h *uploadHandler) handleCreateUseCaseError(w http.ResponseWriter, file *model.File) {
	if err := os.Remove(h.config.Server.Storage.Path + "/" + file.FileId); err != nil {
		slog.Error("Could not remove file from fs", "fileId", file.FileId)
	}

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
