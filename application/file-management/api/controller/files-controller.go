package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"raspstore.github.io/file-manager/api/dto"
	"raspstore.github.io/file-manager/model"
	"raspstore.github.io/file-manager/repository"
	"raspstore.github.io/file-manager/system"
	"raspstore.github.io/file-manager/validators"
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
		send(w, dto.ErrorResponse{Message: "multipart form headers have wrong length", Code: "UP01"})
		return
	}

	header := fileHeader[0]

	aux, err := header.Open()

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		send(w, dto.ErrorResponse{Message: "could not open file", Code: "UP02"})
		return
	}

	filename := header.Filename

	file, err := ioutil.ReadAll(aux)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		send(w, dto.ErrorResponse{Message: "could not read file", Reason: err.Error(), Code: "UP03"})
		return
	}

	fileRef := model.NewFile(filename, r.Header.Get("UID"), uint32(header.Size))

	buffer := bytes.Buffer{}

	if _, err := buffer.Write(file); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		send(w, dto.ErrorResponse{Message: "could not write file to server, try again later", Reason: err.Error(), Code: "UP05"})
		return
	}

	if err := f.ds.Save(fileRef, buffer); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		send(w, dto.ErrorResponse{Message: "could not write file to server, try again later", Reason: err.Error(), Code: "UP06"})
		return
	}

	if err := f.repo.Save(fileRef); err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		// delete file locally
		if delErr := f.ds.Delete(fileRef.Uri); delErr != nil {
			log.Println("WARNING: FILE ", fileRef.Id, " COULD NOT BE DELETED LOCALLY. MANUAL DELETE IS REQUIRED.")
		}

		send(w, dto.ErrorResponse{Message: "could not save file in database. Try again later.", Reason: err.Error(), Code: "UP04"})
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
		send(w, dto.ErrorResponse{Message: "file does not exists", Code: "DW00"})
		return
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		send(w, dto.ErrorResponse{Message: "could not find file, try again later", Reason: err.Error(), Code: "DW01"})
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

	if err := validators.ValidateId(id); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		send(w, dto.ErrorResponse{Message: "Wrong Id", Reason: err.Error(), Code: "DL00"})
		return
	}

	file, err := f.repo.FindById(id)

	if err == mongo.ErrNoDocuments {
		w.WriteHeader(http.StatusNotFound)
		send(w, dto.ErrorResponse{Message: "file does not exists", Code: "DL01"})
		return
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		send(w, dto.ErrorResponse{Message: "could not find file, try again later", Reason: err.Error(), Code: "DL02"})
		return
	}

	if file == nil {
		w.WriteHeader(http.StatusNotFound)
		send(w, dto.ErrorResponse{Message: "file does no exists", Reason: err.Error(), Code: "DL03"})
		return
	}

	if err := f.repo.Delete(id); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		send(w, dto.ErrorResponse{Message: "Could not delete file. Try again later.", Reason: err.Error(), Code: "DL04"})
		return
	}

	if err := f.ds.Delete(file.Uri); err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		// writing file in database again

		f.repo.Save(file)

		send(w, dto.ErrorResponse{Message: "Could not delete file from server. Try again later.", Reason: err.Error(), Code: "DL05"})
		return
	}

}

func (f *filesController) Update(w http.ResponseWriter, r *http.Request) {

	if !validators.ValidateBody(r.Header.Get("Content-Type"), "multipart/form-data") {
		w.WriteHeader(http.StatusUnprocessableEntity)
		send(w, dto.ErrorResponse{Message: "Request body must be a multipart/form-data", Code: "UP00"})
		return
	}

	params := mux.Vars(r)

	id := params["id"]

	if err := validators.ValidateId(id); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		send(w, dto.ErrorResponse{Message: "Wrong Id", Reason: err.Error(), Code: "UP01"})
		return
	}

	fileRef, err := f.repo.FindById(id)

	if err == mongo.ErrNoDocuments {
		w.WriteHeader(http.StatusNotFound)
		send(w, dto.ErrorResponse{Message: fmt.Sprintf("File with id %s not found.", id), Reason: err.Error(), Code: "UP02"})
		return
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		send(w, dto.ErrorResponse{Message: "Could not retrieve file information. Try again later.", Reason: err.Error(), Code: "UP03"})
		return
	}

	r.ParseMultipartForm(1000000)

	datas := r.MultipartForm

	fileHeader := datas.File["file"]

	if len(fileHeader) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		send(w, dto.ErrorResponse{Message: "multipart form headers have wrong length", Code: "UP04"})
		return
	}

	header := fileHeader[0]

	aux, err := header.Open()

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		send(w, dto.ErrorResponse{Message: "could not open file", Reason: err.Error(), Code: "UP05"})
		return
	}

	file, err := ioutil.ReadAll(aux)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		send(w, dto.ErrorResponse{Message: "could not read file", Reason: err.Error(), Code: "UP06"})
		return
	}

	buffer := bytes.Buffer{}

	if _, err := buffer.Write(file); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		send(w, dto.ErrorResponse{Message: "could not write file to server, try again later", Reason: err.Error(), Code: "UP07"})
		return
	}

	fileRef.Filename = header.Filename

	if err := f.ds.Save(fileRef, buffer); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		send(w, dto.ErrorResponse{Message: "could not write file to server, try again later", Reason: err.Error(), Code: "UP08"})
		return
	}

	fileRef.UpdatedBy = r.Header.Get("UID")

	if err := f.repo.Update(fileRef); err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		// delete file locally
		if delErr := f.ds.Delete(fileRef.Uri); delErr != nil {
			log.Println("WARNING: FILE ", fileRef.Id, " COULD NOT BE DELETED LOCALLY. MANUAL DELETE IS REQUIRED.")
		}

		send(w, dto.ErrorResponse{Message: "could not save file in database. Try again later.", Reason: err.Error(), Code: "UP09"})
		return
	}

	send(w, fileRef)

}

func (f *filesController) ListFiles(w http.ResponseWriter, r *http.Request) {
	files, err := f.repo.FindAll()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		send(w, dto.ErrorResponse{Message: "could not list files. Try again later.", Reason: err.Error(), Code: "LF01"})
		return
	}

	send(w, files)
}

func send(w http.ResponseWriter, obj interface{}) {
	w.Header().Set("Content-Type", "application/json")
	jsonResponse, err := json.Marshal(obj)
	if err != nil {
		return
	}
	w.Write(jsonResponse)
}
