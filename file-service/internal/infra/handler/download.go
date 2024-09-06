package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/lestrrat-go/jwx/jwt"
	m "github.com/murilo-bracero/raspstore/file-service/internal/infra/middleware"
	"github.com/murilo-bracero/raspstore/file-service/internal/infra/repository"
	"github.com/murilo-bracero/raspstore/file-service/internal/infra/response"
)

func (h *Handler) Download(w http.ResponseWriter, r *http.Request) {
	fileId := chi.URLParam(r, "fileId")
	usr := r.Context().Value(m.UserClaimsCtxKey).(jwt.Token)
	traceId := r.Context().Value(chiMiddleware.RequestIDKey).(string)

	fileRep, err := h.fileFacade.FindById(usr.Subject(), fileId)

	if err == repository.ErrFileDoesNotExists {
		response.NotFound(w, traceId)
		return
	}

	if err != nil {
		response.InternalServerError(w, traceId)
		return
	}

	file, err := h.downloadUseCase.Execute(r.Context(), fileId)

	if err != nil {
		response.InternalServerError(w, traceId)
		return
	}

	defer file.Close()

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileRep.Filename))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", fileRep.Size))

	http.ServeContent(w, r, fileRep.Filename, time.Now(), file)
}
