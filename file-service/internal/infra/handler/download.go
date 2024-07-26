package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/murilo-bracero/raspstore/file-service/internal/application/facade"
	"github.com/murilo-bracero/raspstore/file-service/internal/application/repository"
	"github.com/murilo-bracero/raspstore/file-service/internal/application/usecase"
	m "github.com/murilo-bracero/raspstore/file-service/internal/infra/middleware"
)

type DownloadHandler interface {
	Download(w http.ResponseWriter, r *http.Request)
}

type downloadHandler struct {
	downloadUseCase usecase.DownloadFileUseCase
	fileFacade      facade.FileFacade
}

func NewDownloadHandler(downloadUseCase usecase.DownloadFileUseCase, fileFacade facade.FileFacade) DownloadHandler {
	return &downloadHandler{downloadUseCase: downloadUseCase, fileFacade: fileFacade}
}

func (h *downloadHandler) Download(w http.ResponseWriter, r *http.Request) {
	fileId := chi.URLParam(r, "fileId")
	usr := r.Context().Value(m.UserClaimsCtxKey).(jwt.Token)

	fileRep, err := h.fileFacade.FindById(usr.Subject(), fileId)

	if err == repository.ErrFileDoesNotExists {
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
