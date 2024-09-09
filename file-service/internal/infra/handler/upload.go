package handler

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/murilo-bracero/raspstore/file-service/internal/domain/entity"
	"github.com/murilo-bracero/raspstore/file-service/internal/domain/model"
	m "github.com/murilo-bracero/raspstore/file-service/internal/infra/middleware"
)

func (h *Handler) Upload(w http.ResponseWriter, r *http.Request) {
	usr := r.Context().Value(m.UserClaimsCtxKey).(jwt.Token)
	traceId := r.Context().Value(middleware.RequestIDKey).(string)
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		slog.Error("Could not allocate MultipartForm parser", "traceId", traceId, "error", err)
		unprocessableEntity(w, traceId)
		return
	}

	file, header, err := r.FormFile("file")

	if err != nil {
		slog.Error("Could not open Multipart Form file", "traceId", traceId, "error", err)
		unprocessableEntity(w, traceId)
		return
	}

	defer file.Close()

	fm := entity.NewFile(header.Filename, header.Size, false, usr.Subject())

	if err := h.fileSystemFacade.Upload(traceId, fm, file); err != nil {
		internalServerError(w, traceId)
		return
	}

	if err := h.createFileUseCase.Execute(fm); err != nil {
		h.handleCreateUseCaseError(w, fm, traceId)
		return
	}

	created(w, &model.UploadSuccessResponse{
		FileId:   fm.FileId,
		Filename: fm.Filename,
		OwnerId:  fm.Owner,
	}, traceId)
}

func (h *Handler) handleCreateUseCaseError(w http.ResponseWriter, file *entity.File, traceId string) {
	if err := os.Remove(h.config.Storage.Path + "/storage/" + file.FileId); err != nil {
		slog.Error("Could not remove file from fs", "fileId", file.FileId)
	}

	internalServerError(w, traceId)
}
