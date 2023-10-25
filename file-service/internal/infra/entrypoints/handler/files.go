package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	u "github.com/murilo-bracero/raspstore/file-service/internal/application/common/utils"
	"github.com/murilo-bracero/raspstore/file-service/internal/application/parser"
	"github.com/murilo-bracero/raspstore/file-service/internal/application/usecase"
	"github.com/murilo-bracero/raspstore/file-service/internal/domain/entity"
	e "github.com/murilo-bracero/raspstore/file-service/internal/domain/exceptions"
	"github.com/murilo-bracero/raspstore/file-service/internal/domain/model/request"
	"github.com/murilo-bracero/raspstore/file-service/internal/infra/validators"
)

type FilesHandler interface {
	ListFiles(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type filesHandler struct {
	listUseCase   usecase.ListFilesUseCase
	updateUseCase usecase.UpdateFileUseCase
	deleteUseCase usecase.DeleteFileUseCase
}

func NewFilesHandler(listUseCase usecase.ListFilesUseCase, updateUseCase usecase.UpdateFileUseCase, deleteUseCase usecase.DeleteFileUseCase) FilesHandler {
	return &filesHandler{listUseCase: listUseCase, updateUseCase: updateUseCase, deleteUseCase: deleteUseCase}
}

func (f *filesHandler) ListFiles(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	size, _ := strconv.Atoi(r.URL.Query().Get("size"))
	filename := r.URL.Query().Get("filename")
	secretQuery := r.URL.Query().Get("secret")

	secret, _ := strconv.ParseBool(secretQuery)

	filesPage, err := f.listUseCase.Execute(r.Context(), page, size, filename, secret)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	nextUrl := buildNextUrl(filesPage, r.Host, page, size)

	u.Send(w, parser.FilePageResponseParser(page, size, filesPage, nextUrl))
}

func (f *filesHandler) Update(w http.ResponseWriter, r *http.Request) {
	traceId := r.Context().Value(chiMiddleware.RequestIDKey).(string)

	var req request.UpdateFileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
		return
	}

	if err := validators.ValidateUpdateFileRequest(&req); err != nil {
		u.HandleBadRequest(w, traceId, "ERR001", err.Error())
		return
	}

	fileId := chi.URLParam(r, "id")

	file := &entity.File{
		FileId:   fileId,
		Secret:   req.Secret,
		Filename: req.Filename,
		Editors:  req.Editors,
		Viewers:  req.Viewers,
	}

	fileMetadata, err := f.updateUseCase.Execute(r.Context(), file)

	if err == e.ErrFileDoesNotExists {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	u.Send(w, fileMetadata)
}

func (f *filesHandler) Delete(w http.ResponseWriter, r *http.Request) {
	fileId := chi.URLParam(r, "id")

	if err := f.deleteUseCase.Execute(r.Context(), fileId); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusNoContent)
}

func buildNextUrl(filesPage *entity.FilePage, host string, page int, size int) (nextUrl string) {
	if len(filesPage.Content) == size {
		nextUrl = fmt.Sprintf("%s/file-service/files?page=%d&size=%d", host, page+1, size)
	}

	return
}
