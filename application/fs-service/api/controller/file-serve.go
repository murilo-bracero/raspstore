package controller

import (
	"io"
	"log"
	"net/http"
	"os"

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

	fileBytes, err := io.ReadAll(file)

	if err != nil {
		traceId := r.Context().Value(middleware.RequestIDKey).(string)
		log.Printf("[ERROR] - [%s]: Could not read file buffer due to error: %s", traceId, err.Error())
		InternalServerError(w, traceId)
		return
	}

	filerep.Write(fileBytes)
}

func (f *fileServeController) Download(w http.ResponseWriter, r *http.Request) {
}
