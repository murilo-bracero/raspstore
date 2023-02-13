package controller

import "net/http"

type FileServeController interface {
	Upload(w http.ResponseWriter, r *http.Request)
	Download(w http.ResponseWriter, r *http.Request)
}

type fileServeController struct {
}

func NewFileServeController() FileServeController {
	return &fileServeController{}
}

func (f *fileServeController) Upload(w http.ResponseWriter, r *http.Request) {

}

func (f *fileServeController) Download(w http.ResponseWriter, r *http.Request) {

}
