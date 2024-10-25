package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/murilo-bracero/raspstore/file-service/internal/infra/repository"
)

func (h *Handler) Download(w http.ResponseWriter, r *http.Request) {
	fileId := chi.URLParam(r, "fileId")
	usr := r.Context().Value(UserClaimsCtxKey).(jwt.Token)
	traceId := r.Context().Value(chiMiddleware.RequestIDKey).(string)

	fileRep, err := h.fileFacade.FindById(usr.Subject(), fileId)

	if err == repository.ErrFileDoesNotExists {
		notFound(w, traceId)
		return
	}

	if err != nil {
		internalServerError(w, traceId)
		return
	}

	file, err := h.fileSystemFacade.Download(traceId, usr.Subject(), fileRep.FileId)

	if err != nil {
		internalServerError(w, traceId)
		return
	}

	defer file.Close()

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileRep.Filename))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", fileRep.Size))

	http.ServeContent(w, r, fileRep.Filename, time.Now(), file)
}
