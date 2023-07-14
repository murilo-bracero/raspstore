package handler

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	rMiddleware "github.com/murilo-bracero/raspstore/commons/pkg/middleware"
	"github.com/murilo-bracero/raspstore/file-info-service/proto/v1/file-info-service/pb"
	v1 "raspstore.github.io/fs-service/api/v1"
	u "raspstore.github.io/fs-service/internal/api/utils"
	"raspstore.github.io/fs-service/internal/usecase"
)

type FileServeHandler interface {
	Upload(w http.ResponseWriter, r *http.Request)
	Download(w http.ResponseWriter, r *http.Request)
}

type fileServeHandler struct {
	uploadUseCase   usecase.UploadFileUseCase
	downloadUseCase usecase.DownloadFileUseCase
}

func NewFileServeHandler(uuc usecase.UploadFileUseCase, duc usecase.DownloadFileUseCase) FileServeHandler {
	return &fileServeHandler{uploadUseCase: uuc, downloadUseCase: duc}
}

func (f *fileServeHandler) Upload(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(rMiddleware.UserIdKey).(string)
	traceId := r.Context().Value(middleware.RequestIDKey).(string)
	r.ParseMultipartForm(32 << 20)

	file, header, err := r.FormFile("file")
	path := r.FormValue("path")

	if path == "" {
		u.HandleBadRequest(w, traceId, "File filed is null or malformed", "ERR001")
		return
	}

	if err != nil {
		log.Printf("[ERROR] - [%s]: Could not open Multipart Form file: %s", traceId, err.Error())
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
		return
	}

	defer file.Close()

	req := &pb.CreateFileMetadataRequest{
		OwnerId:  userId,
		Filename: header.Filename,
		Size:     header.Size,
		Path:     path,
	}

	res, err := f.uploadUseCase.Execute(r.Context(), req, file)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	u.Created(w, &v1.UploadSuccessResponse{
		FileId:   res.FileId,
		Filename: res.Filename,
		Path:     res.Path,
		OwnerId:  res.OwnerId,
	})
}

func (f *fileServeHandler) Download(w http.ResponseWriter, r *http.Request) {

	fileId := chi.URLParam(r, "fileId")

	downloadRep, err := f.downloadUseCase.Execute(r.Context(), fileId)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	defer downloadRep.File.Close()

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", downloadRep.Filename))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", downloadRep.FileSize))

	http.ServeContent(w, r, downloadRep.Filename, time.Now(), downloadRep.File)
}
