package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/murilo-bracero/raspstore/file-service/internal/domain/entity"
	"github.com/murilo-bracero/raspstore/file-service/internal/domain/model"
	"github.com/murilo-bracero/raspstore/file-service/internal/infra/mapper"
	m "github.com/murilo-bracero/raspstore/file-service/internal/infra/middleware"
	"github.com/murilo-bracero/raspstore/file-service/internal/infra/repository"
	"github.com/murilo-bracero/raspstore/file-service/internal/infra/response"
	"github.com/murilo-bracero/raspstore/file-service/internal/infra/validator"
)

func (f *Handler) ListFiles(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	size, _ := strconv.Atoi(r.URL.Query().Get("size"))

	filename := r.URL.Query().Get("filename")
	secretQuery := r.URL.Query().Get("secret")

	traceId := r.Context().Value(chiMiddleware.RequestIDKey).(string)
	user := r.Context().Value(m.UserClaimsCtxKey).(jwt.Token)

	secret, _ := strconv.ParseBool(secretQuery)

	filesPage, err := f.fileFacade.FindAll(traceId, user.Subject(), page, size, filename, secret)

	if err != nil {
		response.InternalServerError(w, traceId)
		return
	}

	response.Ok(w, mapper.MapFilePageResponse(page, size, filesPage, r.Host), traceId)
}

func (f *Handler) FindById(w http.ResponseWriter, r *http.Request) {
	traceId := r.Context().Value(chiMiddleware.RequestIDKey).(string)
	user := r.Context().Value(m.UserClaimsCtxKey).(jwt.Token)

	fileId := chi.URLParam(r, "id")

	entity, err := f.fileFacade.FindById(user.Subject(), fileId)

	if err == repository.ErrFileDoesNotExists {
		response.NotFound(w, traceId)
		return
	}

	response.Ok(w, entity, traceId)
}

func (f *Handler) Update(w http.ResponseWriter, r *http.Request) {
	traceId := r.Context().Value(chiMiddleware.RequestIDKey).(string)

	var req model.UpdateFileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.UnprocessableEntity(w, traceId)
		return
	}

	if err := validator.ValidateUpdateFileRequest(&req); err != nil {
		response.BadRequest(w, model.ErrorResponse{
			Message: err.Error(),
		}, traceId)
		return
	}

	fileId := chi.URLParam(r, "id")

	file := &entity.File{
		FileId:   fileId,
		Secret:   req.Secret,
		Filename: req.Filename,
	}

	fileMetadata, err := f.updateFileUseCase.Execute(r.Context(), file)

	if err == repository.ErrFileDoesNotExists {
		response.NotFound(w, traceId)
		return
	}

	if err != nil {
		response.InternalServerError(w, traceId)
		return
	}

	response.Ok(w, fileMetadata, traceId)
}

func (f *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	fileId := chi.URLParam(r, "id")

	slog.Info(fileId)

	traceId := r.Context().Value(chiMiddleware.RequestIDKey).(string)
	user := r.Context().Value(m.UserClaimsCtxKey).(jwt.Token)

	if err := f.fileFacade.DeleteById(traceId, user.Subject(), fileId); err != nil {
		response.InternalServerError(w, traceId)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
