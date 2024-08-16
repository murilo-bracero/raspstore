package handler

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/murilo-bracero/raspstore/file-service/internal/application/usecase"
	"github.com/murilo-bracero/raspstore/file-service/internal/domain/entity"
	"github.com/murilo-bracero/raspstore/file-service/internal/domain/model"
	"github.com/murilo-bracero/raspstore/file-service/internal/infra/config"
	m "github.com/murilo-bracero/raspstore/file-service/internal/infra/middleware"
	"github.com/murilo-bracero/raspstore/file-service/internal/infra/response"
)

type UploadHandler interface {
	Upload(w http.ResponseWriter, r *http.Request)
}

type uploadHandler struct {
	config            *config.Config
	uploadUseCase     usecase.UploadFileUseCase
	createFileUseCase usecase.CreateFileUseCase
}

func NewUploadHandler(config *config.Config, uploadUseCase usecase.UploadFileUseCase, createFileUseCase usecase.CreateFileUseCase) UploadHandler {
	return &uploadHandler{config: config, uploadUseCase: uploadUseCase, createFileUseCase: createFileUseCase}
}

func (h *uploadHandler) Upload(w http.ResponseWriter, r *http.Request) {
	usr := r.Context().Value(m.UserClaimsCtxKey).(jwt.Token)
	traceId := r.Context().Value(middleware.RequestIDKey).(string)
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		slog.Error("Could not allocate MultipartForm parser", "traceId", traceId, "error", err)
		response.UnprocessableEntity(w, traceId)
		return
	}

	file, header, err := r.FormFile("file")

	if err != nil {
		slog.Error("Could not open Multipart Form file", "traceId", traceId, "error", err)
		response.UnprocessableEntity(w, traceId)
		return
	}

	defer file.Close()

	fm := entity.NewFile(header.Filename, header.Size, false, usr.Subject())

	if err := h.uploadUseCase.Execute(r.Context(), fm, file); err != nil {
		response.InternalServerError(w, traceId)
		return
	}

	if err := h.createFileUseCase.Execute(fm); err != nil {
		h.handleCreateUseCaseError(w, fm, traceId)
		return
	}

	response.Created(w, &model.UploadSuccessResponse{
		FileId:   fm.FileId,
		Filename: fm.Filename,
		OwnerId:  fm.Owner,
	}, traceId)
}

func (h *uploadHandler) handleCreateUseCaseError(w http.ResponseWriter, file *entity.File, traceId string) {
	if err := os.Remove(h.config.Storage.Path + "/storage/" + file.FileId); err != nil {
		slog.Error("Could not remove file from fs", "fileId", file.FileId)
	}

	response.InternalServerError(w, traceId)
}
