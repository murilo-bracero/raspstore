package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	rMiddleware "github.com/murilo-bracero/raspstore/commons/pkg/middleware"
	"go.mongodb.org/mongo-driver/mongo"
	v1 "raspstore.github.io/file-manager/api/v1"
	"raspstore.github.io/file-manager/internal/repository"
)

const maxListSize = 50

type FilesHandler interface {
	ListFiles(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type filesHandler struct {
	repo repository.FilesRepository
}

func NewFilesHandler(repo repository.FilesRepository) FilesHandler {
	return &filesHandler{repo: repo}
}

func (f *filesHandler) ListFiles(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	size, _ := strconv.Atoi(r.URL.Query().Get("size"))

	if size == 0 || size > maxListSize {
		size = maxListSize
	}

	userId := r.Context().Value(rMiddleware.UserIdKey).(string)

	filesPage, err := f.repo.FindAll(userId, page, size)

	if err != nil {
		traceId := r.Context().Value(chiMiddleware.RequestIDKey).(string)
		log.Printf("[ERROR] - [%s]: Could list files due to error: %s", traceId, err.Error())
		v1.InternalServerError(w, traceId)
		return
	}

	nextUrl := ""

	if len(filesPage.Content) == size {
		nextUrl = fmt.Sprintf("%s/file-info-service/files?page=%d&size=%d", r.Host, page+1, size)
	}

	v1.Send(w, v1.FileResponse{
		Page:          page,
		Size:          size,
		TotalElements: filesPage.Count,
		Next:          nextUrl,
		Content:       filesPage.Content,
	})
}

func (f *filesHandler) Update(w http.ResponseWriter, r *http.Request) {
	var req v1.UpdateFileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
		return
	}

	fileId := chi.URLParam(r, "id")
	userId := r.Context().Value(rMiddleware.UserIdKey).(string)
	traceId := r.Context().Value(chiMiddleware.RequestIDKey).(string)

	file, err := f.repo.FindById(userId, fileId)

	if err == mongo.ErrNoDocuments {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if err != nil {
		log.Printf("[ERROR] - [%s]: Could not search file with id %s in database: %s", traceId, fileId, err.Error())
		v1.InternalServerError(w, traceId)
		return
	}

	log.Printf("[INFO] - [%s]: File with id=%s found", traceId, fileId)

	if req.Filename != "" {
		file.Filename = req.Filename
	}

	if req.Path != "" {
		file.Path = req.Path
	}

	if req.Viewers != nil || len(req.Viewers) != 0 {
		file.Viewers = req.Viewers
	}

	if req.Editors != nil || len(req.Editors) != 0 {
		file.Editors = req.Editors
	}

	if err := f.repo.Update(userId, file); err != nil {
		log.Printf("[ERROR] - [%s]: Could not update file with id %s in database: %s", traceId, fileId, err.Error())
		v1.InternalServerError(w, traceId)
		return
	}

	log.Printf("[INFO] - [%s]: File with id=%s updated successfully", traceId, fileId)

	fileMetadata, err := f.repo.FindByIdLookup(userId, fileId)

	if err != nil {
		log.Printf("[ERROR] - [%s]: Could not search lookup file with id %s in database: %s", traceId, fileId, err.Error())
		v1.InternalServerError(w, traceId)
		return
	}

	v1.Send(w, fileMetadata)
}

func (f *filesHandler) Delete(w http.ResponseWriter, r *http.Request) {
	fileId := chi.URLParam(r, "id")

	userId := r.Context().Value(rMiddleware.UserIdKey).(string)

	if err := f.repo.Delete(userId, fileId); err != nil {
		traceId := r.Context().Value(chiMiddleware.RequestIDKey).(string)
		log.Printf("[ERROR] - [%s]: Could not delete file with id %s in database: %s", traceId, fileId, err.Error())
		v1.InternalServerError(w, traceId)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
