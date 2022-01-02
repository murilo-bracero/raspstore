package controller

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"raspstore.github.io/file-manager/api/dto"
	"raspstore.github.io/file-manager/model"
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

func (f *filesController) Upload(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(1000000)

	datas := r.MultipartForm

	fileHeader := datas.File["file"]

	if len(fileHeader) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		er := new(dto.ErrorResponse)
		er.Message = "multipart form headers have wrong length"
		er.Code = "UP01"
		send(w, er)
		return
	}

	header := fileHeader[0]

	aux, err := header.Open()

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		er := new(dto.ErrorResponse)
		er.Message = "could not open file: uploaded file is corrupted or network connectivity is bad, try again later"
		er.Code = "UP02"
		send(w, er)
		return
	}

	filename := header.Filename

	file, err := ioutil.ReadAll(aux)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		er := new(dto.ErrorResponse)
		er.Message = "could not read file: uploaded file is corrupted or network connectivity is bad, try again later"
		er.Reason = err.Error()
		er.Code = "UP03"
		send(w, er)
		return
	}

	fileRef := model.NewFile(filename, r.Header.Get("UID"), uint32(header.Size))

	buffer := bytes.Buffer{}

	if _, err := buffer.Write(file); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		er := new(dto.ErrorResponse)
		er.Message = "could not write file to server, try again later"
		er.Reason = err.Error()
		er.Code = "UP05"
		send(w, er)
		return
	}

	if err := f.ds.Save(fileRef, buffer); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		er := new(dto.ErrorResponse)
		er.Message = "could not write file to server, try again later"
		er.Reason = err.Error()
		er.Code = "UP06"
		send(w, er)
		return
	}

	if err := f.repo.Save(fileRef); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		er := new(dto.ErrorResponse)
		er.Message = "could not save file in database. Try again later."
		er.Reason = err.Error()
		er.Code = "UP04"

		// delete file locally
		if delErr := f.ds.Delete(fileRef.Uri); delErr != nil {
			log.Println("WARNING: FILE ", fileRef.Id, " COULD NOT BE DELETED LOCALLY. MANUAL DELETE IS REQUIRED.")
		}

		send(w, er)
		return
	}

	send(w, fileRef)

}

func (f *filesController) Download(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id := params["id"]

	file, err := f.repo.FindById(id)

	if err == mongo.ErrNoDocuments {
		w.WriteHeader(http.StatusNotFound)
		er := new(dto.ErrorResponse)
		er.Message = "file does not exists"
		er.Code = "DW00"
		send(w, er)
		return
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		er := new(dto.ErrorResponse)
		er.Message = "could not find file, try again later"
		er.Reason = err.Error()
		er.Code = "DW01"
		send(w, er)
		return
	}

	if r.URL.Query().Get("download") == "true" {
		w.Header().Set("Content-Disposition", "attachment; filename="+file.Filename)
		http.ServeFile(w, r, file.Uri)
	} else {
		send(w, file)
	}

}

func (f *filesController) Delete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id := params["id"]

	file, err := f.repo.FindById(id)

	if err == mongo.ErrNoDocuments {
		w.WriteHeader(http.StatusNotFound)
		er := new(dto.ErrorResponse)
		er.Message = "file does not exists"
		er.Code = "DL00"
		send(w, er)
		return
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		er := new(dto.ErrorResponse)
		er.Message = "could not find file, try again later"
		er.Reason = err.Error()
		er.Code = "DL01"
		send(w, er)
		return
	}

	if file == nil {
		w.WriteHeader(http.StatusNotFound)
		er := new(dto.ErrorResponse)
		er.Message = "file does no exists"
		er.Reason = err.Error()
		er.Code = "DL02"
		send(w, er)
		return
	}

	if err := f.repo.Delete(id); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		er := new(dto.ErrorResponse)
		er.Message = "could not delete file. Try again later."
		er.Reason = err.Error()
		er.Code = "UP06"
		send(w, er)
		return
	}

	if err := f.ds.Delete(file.Uri); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		er := new(dto.ErrorResponse)
		er.Message = "could not delete file from server. Try again later."
		er.Reason = err.Error()
		er.Code = "DL03"

		// writing file in database again

		f.repo.Save(file)

		send(w, er)
		return
	}

}

func (f *filesController) Update(w http.ResponseWriter, r *http.Request) {

	if !validateBody(r.Header.Get("Content-Type"), "multipart/form-data") {
		w.WriteHeader(http.StatusUnprocessableEntity)
		er := new(dto.ErrorResponse)
		er.Message = "Request body must be a multipart/form-data"
		er.Code = "UP00"
		send(w, er)
		return
	}

	params := mux.Vars(r)

	id := params["id"]

	fileRef, err := f.repo.FindById(id)

	if err == mongo.ErrNoDocuments {
		w.WriteHeader(http.StatusNotFound)
		er := new(dto.ErrorResponse)
		er.Message = "File with id " + id + " not found."
		er.Reason = err.Error()
		er.Code = "UP01"
		send(w, er)
		return
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		er := new(dto.ErrorResponse)
		er.Message = "Could not retrieve file information. Try again later."
		er.Reason = err.Error()
		er.Code = "UP02"
		send(w, er)
		return
	}

	r.ParseMultipartForm(1000000)

	datas := r.MultipartForm

	fileHeader := datas.File["file"]

	if len(fileHeader) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		er := new(dto.ErrorResponse)
		er.Message = "multipart form headers have wrong length"
		er.Code = "UP03"
		send(w, er)
		return
	}

	header := fileHeader[0]

	aux, err := header.Open()

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		er := new(dto.ErrorResponse)
		er.Message = "could not open file: uploaded file is corrupted or network connectivity is bad, try again later"
		er.Code = "UP04"
		send(w, er)
		return
	}

	file, err := ioutil.ReadAll(aux)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		er := new(dto.ErrorResponse)
		er.Message = "could not read file: uploaded file is corrupted or network connectivity is bad, try again later"
		er.Reason = err.Error()
		er.Code = "UP05"
		send(w, er)
		return
	}

	buffer := bytes.Buffer{}

	if _, err := buffer.Write(file); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		er := new(dto.ErrorResponse)
		er.Message = "could not write file to server, try again later"
		er.Reason = err.Error()
		er.Code = "UP06"
		send(w, er)
		return
	}

	fileRef.Filename = header.Filename

	if err := f.ds.Save(fileRef, buffer); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		er := new(dto.ErrorResponse)
		er.Message = "could not write file to server, try again later"
		er.Reason = err.Error()
		er.Code = "UP07"
		send(w, er)
		return
	}

	fileRef.UpdatedBy = r.Header.Get("UID")

	if err := f.repo.Update(fileRef); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		er := new(dto.ErrorResponse)
		er.Message = "could not save file in database. Try again later."
		er.Reason = err.Error()
		er.Code = "UP08"

		// delete file locally
		if delErr := f.ds.Delete(fileRef.Uri); delErr != nil {
			log.Println("WARNING: FILE ", fileRef.Id, " COULD NOT BE DELETED LOCALLY. MANUAL DELETE IS REQUIRED.")
		}

		send(w, er)
		return
	}

	send(w, fileRef)

}

func (f *filesController) ListFiles(w http.ResponseWriter, r *http.Request) {
	files, err := f.repo.FindAll()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		er := new(dto.ErrorResponse)
		er.Message = "could not list files. Try again later."
		er.Reason = err.Error()
		er.Code = "LF01"
		send(w, er)
		return
	}

	send(w, files)
}

func validateBody(received string, desired string) bool {
	clean := strings.Split(received, ";")

	return clean[0] == desired
}

func send(w http.ResponseWriter, obj interface{}) {
	w.Header().Set("Content-Type", "application/json")
	jsonResponse, err := json.Marshal(obj)
	if err != nil {
		return
	}
	w.Write(jsonResponse)
}
