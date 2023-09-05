package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/murilo-bracero/raspstore/file-service/internal"
	m "github.com/murilo-bracero/raspstore/file-service/internal/api/middleware"
	"github.com/murilo-bracero/raspstore/file-service/internal/usecase"
)

type DownloadHandler interface {
	Download(w http.ResponseWriter, r *http.Request)
}

type downloadHandler struct {
	downloadUseCase usecase.DownloadFileUseCase
	getFileUseCase  usecase.GetFileUseCase
}

func NewDownloadHandler(downloadUseCase usecase.DownloadFileUseCase, getFileUseCase usecase.GetFileUseCase) DownloadHandler {
	return &downloadHandler{downloadUseCase: downloadUseCase, getFileUseCase: getFileUseCase}
}

func (h *downloadHandler) Download(w http.ResponseWriter, r *http.Request) {
	fileId := chi.URLParam(r, "fileId")
	userId := r.Context().Value(m.UserClaimsCtxKey).(string)

	fileRep, err := h.getFileUseCase.Execute(userId, fileId)

	if err == internal.ErrFileDoesNotExists {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	file, err := h.downloadUseCase.Execute(r.Context(), fileId)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	defer file.Close()

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileRep.Filename))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", fileRep.Size))

	http.ServeContent(w, r, fileRep.Filename, time.Now(), file)
}
