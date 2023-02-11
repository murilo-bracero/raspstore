package controller

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5/middleware"
	"raspstore.github.io/file-manager/api/dto"
	"raspstore.github.io/file-manager/repository"
)

const maxListSize = 50

type FilesController interface {
	ListFiles(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type filesController struct {
	repo repository.FilesRepository
}

func NewFilesController(repo repository.FilesRepository) FilesController {
	return &filesController{repo: repo}
}

func (f *filesController) ListFiles(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	size, _ := strconv.Atoi(r.URL.Query().Get("size"))

	if size == 0 || size > maxListSize {
		size = maxListSize
	}

	filesPage, err := f.repo.FindAll(page, size)

	if err != nil {
		traceId := r.Context().Value(middleware.RequestIDKey).(string)
		log.Printf("[ERROR] - [%s]: Could list users due to error: %s", traceId, err.Error())
		InternalServerError(w, traceId)
		return
	}

	nextUrl := ""

	if len(filesPage.Content) == size {
		nextUrl = fmt.Sprintf("%s/users-service/users?page=%d&size=%d", r.Host, page+1, size)
	}

	Send(w, dto.FileMetadataList{
		Page:          page,
		Size:          size,
		TotalElements: filesPage.Count,
		Next:          nextUrl,
		Content:       filesPage.Content,
	})
}

func (f *filesController) Update(w http.ResponseWriter, r *http.Request) {
}

func (f *filesController) Delete(w http.ResponseWriter, r *http.Request) {

}
