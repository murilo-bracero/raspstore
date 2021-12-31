package controller

import (
	"net/http"

	"raspstore.github.io/file-manager/repository"
	"raspstore.github.io/file-manager/system"
)

type FilesController interface {
	Upload(w http.ResponseWriter, r *http.Request)
	Download(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	ListFiles(w http.ResponseWriter, r *http.Request)
}

type filesController struct {
	repo repository.FilesRepository
	ds   system.DiskStore
}

func NewFilesController(repo repository.FilesRepository, ds system.DiskStore) FilesController {
	return &filesController{repo: repo, ds: ds}
}

func (f *filesController) Upload(w http.ResponseWriter, r *http.Request)    {}
func (f *filesController) Download(w http.ResponseWriter, r *http.Request)  {}
func (f *filesController) Delete(w http.ResponseWriter, r *http.Request)    {}
func (f *filesController) Update(w http.ResponseWriter, r *http.Request)    {}
func (f *filesController) ListFiles(w http.ResponseWriter, r *http.Request) {}
