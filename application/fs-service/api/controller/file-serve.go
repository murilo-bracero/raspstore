package controller

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/murilo-bracero/raspstore-protofiles/file-info-service/pb"
	"raspstore.github.io/fs-service/api/dto"
	"raspstore.github.io/fs-service/internal"
	"raspstore.github.io/fs-service/usecase"
)

type FileServeController interface {
	Upload(w http.ResponseWriter, r *http.Request)
	Download(w http.ResponseWriter, r *http.Request)
}

type fileServeController struct {
	fiuc usecase.FileInfoUseCase
}

func NewFileServeController(fiuc usecase.FileInfoUseCase) FileServeController {
	return &fileServeController{fiuc: fiuc}
}

func (f *fileServeController) Upload(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)

	file, header, err := r.FormFile("file")

	if err != nil {
		traceId := r.Context().Value(middleware.RequestIDKey).(string)
		log.Printf("[ERROR] - [%s]: Could list files due to error: %s", traceId, err.Error())
		BadRequest(w, dto.ErrorResponse{
			TraceId: traceId,
			Message: "File filed is null or malformed",
			Code:    "ERR001",
		})
		return
	}

	defer file.Close()

	req := &pb.CreateFileMetadataRequest{
		//TODO:Get user ID from JWT
		OwnerId:  "e9e28c79-a5e8-4545-bd32-e536e690bd4a",
		Filename: header.Filename,
		Size:     header.Size,
		Path:     r.FormValue("path"),
	}

	res, err := f.fiuc.CreateFileMetadata(req)

	if err != nil {
		traceId := r.Context().Value(middleware.RequestIDKey).(string)
		log.Printf("[ERROR] - [%s]: Could not create file in file-info-service due to error: %s", traceId, err.Error())
		InternalServerError(w, traceId)
		return
	}

	filerep, err := os.Create(internal.StoragePath() + "/" + res.FileId)

	if err != nil {
		traceId := r.Context().Value(middleware.RequestIDKey).(string)
		log.Printf("[ERROR] - [%s]: Could not create file in fs due to error: %s", traceId, err.Error())
		InternalServerError(w, traceId)
		return
	}

	defer filerep.Close()

	_, err = io.Copy(filerep, file)

	if err != nil {
		traceId := r.Context().Value(middleware.RequestIDKey).(string)
		log.Printf("[ERROR] - [%s]: Could not read file buffer due to error: %s", traceId, err.Error())
		InternalServerError(w, traceId)
		return
	}
}

func (f *fileServeController) Download(w http.ResponseWriter, r *http.Request) {

	fileId := chi.URLParam(r, "fileId")
	//TODO:Get user ID from JWT
	userId := "e9e28c79-a5e8-4545-bd32-e536e690bd4a"

	fileInfo, err := f.fiuc.GetFileMetadataById(fileId)

	if err != nil {
		traceId := r.Context().Value(middleware.RequestIDKey).(string)
		log.Printf("[ERROR] - [%s]: Could not retrieve file from file-info-service due to error: %s", traceId, err.Error())
		InternalServerError(w, traceId)
		return
	}

	if !f.userHasPermission(fileInfo, userId) {
		http.Error(w, "Not Found", http.StatusUnauthorized)
		return
	}

	if err != nil {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	file, err := os.Open(internal.StoragePath() + "/" + fileInfo.FileId)

	if err != nil {
		traceId := r.Context().Value(middleware.RequestIDKey).(string)
		log.Printf("[ERROR] - [%s]: Could not open file in fs due to error: %s", traceId, err.Error())
		InternalServerError(w, traceId)
		return
	}

	defer file.Close()

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileInfo.Filename))

	bytes, err := ioutil.ReadAll(file)

	if err != nil {
		traceId := r.Context().Value(middleware.RequestIDKey).(string)
		log.Printf("[ERROR] - [%s]: Could not read file in fs due to error: %s", traceId, err.Error())
		InternalServerError(w, traceId)
		return
	}

	_, err = w.Write(bytes)

	if err != nil {
		traceId := r.Context().Value(middleware.RequestIDKey).(string)
		log.Printf("[ERROR] - [%s]: Could not write file to response due to error: %s", traceId, err.Error())
		InternalServerError(w, traceId)
		return
	}
}

func (f *fileServeController) userHasPermission(fileInfo *pb.FileMetadata, userId string) bool {
	return userId == fileInfo.OwnerId
}
