package handler

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5/middleware"
	v1 "github.com/murilo-bracero/raspstore/file-service/api/v1"
	"github.com/murilo-bracero/raspstore/file-service/internal"
	m "github.com/murilo-bracero/raspstore/file-service/internal/api/middleware"
	u "github.com/murilo-bracero/raspstore/file-service/internal/api/utils"
	"github.com/murilo-bracero/raspstore/file-service/internal/converter"
	"github.com/murilo-bracero/raspstore/file-service/internal/model"
	"github.com/murilo-bracero/raspstore/file-service/internal/usecase"
)

type UploadHandler interface {
	Upload(w http.ResponseWriter, r *http.Request)
}

type uploadHandler struct {
	uploadUseCase     usecase.UploadFileUseCase
	createFileUseCase usecase.CreateFileUseCase
}

func NewUploadHandler(uploadUseCase usecase.UploadFileUseCase, createFileUseCase usecase.CreateFileUseCase) UploadHandler {
	return &uploadHandler{uploadUseCase: uploadUseCase, createFileUseCase: createFileUseCase}
}

func (h *uploadHandler) Upload(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(m.UserClaimsCtxKey).(string)
	traceId := r.Context().Value(middleware.RequestIDKey).(string)
	r.ParseMultipartForm(32 << 20)

	file, header, err := r.FormFile("file")

	if err != nil {
		slog.Error("[%s]: Could not open Multipart Form file: %s", traceId, err.Error())
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
		return
	}

	defer file.Close()

	fm := converter.CreateFile(header.Filename, header.Size, false, userId)

	if err := h.uploadUseCase.Execute(r.Context(), fm, file); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if err := h.createFileUseCase.Execute(fm); err != nil {
		handleCreateUseCaseError(w, fm)
		return
	}

	u.Created(w, &v1.UploadSuccessResponse{
		FileId:   fm.FileId,
		Filename: fm.Filename,
		OwnerId:  fm.Owner,
	})
}

func handleCreateUseCaseError(w http.ResponseWriter, file *model.File) {
	if err := os.Remove(internal.StoragePath() + "/" + file.FileId); err != nil {
		slog.Error("Could not remove file from fs, need to do it manually: fileId=%s", file.FileId)
	}

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
