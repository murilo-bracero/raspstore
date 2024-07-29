package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/murilo-bracero/raspstore/file-service/internal/application/facade"
	"github.com/murilo-bracero/raspstore/file-service/internal/application/repository"
	"github.com/murilo-bracero/raspstore/file-service/internal/application/usecase"
	"github.com/murilo-bracero/raspstore/file-service/internal/domain/entity"
	"github.com/murilo-bracero/raspstore/file-service/internal/domain/model"
	"github.com/murilo-bracero/raspstore/file-service/internal/infra/mapper"
	m "github.com/murilo-bracero/raspstore/file-service/internal/infra/middleware"
	u "github.com/murilo-bracero/raspstore/file-service/internal/infra/utils"
	"github.com/murilo-bracero/raspstore/file-service/internal/infra/validator"
)

type FilesHandler interface {
	ListFiles(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type filesHandler struct {
	fileFacade    facade.FileFacade
	updateUseCase usecase.UpdateFileUseCase
}

func NewFilesHandler(fileFacade facade.FileFacade, updateUseCase usecase.UpdateFileUseCase) FilesHandler {
	return &filesHandler{fileFacade: fileFacade, updateUseCase: updateUseCase}
}

func (f *filesHandler) ListFiles(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	size, _ := strconv.Atoi(r.URL.Query().Get("size"))

	filename := r.URL.Query().Get("filename")
	secretQuery := r.URL.Query().Get("secret")

	traceId := r.Context().Value(chiMiddleware.RequestIDKey).(string)
	user := r.Context().Value(m.UserClaimsCtxKey).(jwt.Token)

	secret, _ := strconv.ParseBool(secretQuery)

	filesPage, err := f.fileFacade.FindAll(traceId, user.Subject(), page, size, filename, secret)

	if err != nil {
		u.InternalServerError(w, traceId)
		return
	}

	u.Ok(w, mapper.MapFilePageResponse(page, size, filesPage, r.Host), traceId)
}

func (f *filesHandler) Update(w http.ResponseWriter, r *http.Request) {
	traceId := r.Context().Value(chiMiddleware.RequestIDKey).(string)

	var req model.UpdateFileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		u.UnprocessableEntity(w, traceId)
		return
	}

	if err := validator.ValidateUpdateFileRequest(&req); err != nil {
		u.HandleBadRequest(w, traceId, err.Error())
		return
	}

	fileId := chi.URLParam(r, "id")

	file := &entity.File{
		FileId:   fileId,
		Secret:   req.Secret,
		Filename: req.Filename,
	}

	fileMetadata, err := f.updateUseCase.Execute(r.Context(), file)

	if err == repository.ErrFileDoesNotExists {
		u.NotFound(w, traceId)
		return
	}

	if err != nil {
		u.InternalServerError(w, traceId)
		return
	}

	u.Ok(w, fileMetadata, traceId)
}

func (f *filesHandler) Delete(w http.ResponseWriter, r *http.Request) {
	fileId := chi.URLParam(r, "id")

	slog.Info(fileId)

	traceId := r.Context().Value(chiMiddleware.RequestIDKey).(string)
	user := r.Context().Value(m.UserClaimsCtxKey).(jwt.Token)

	if err := f.fileFacade.DeleteById(traceId, user.Subject(), fileId); err != nil {
		u.InternalServerError(w, traceId)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
