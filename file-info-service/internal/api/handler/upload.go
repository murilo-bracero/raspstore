package handler

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	rMiddleware "github.com/murilo-bracero/raspstore/commons/pkg/middleware"
	v1 "github.com/murilo-bracero/raspstore/file-info-service/api/v1"
	u "github.com/murilo-bracero/raspstore/file-info-service/internal/api/utils"
	"github.com/murilo-bracero/raspstore/file-info-service/internal/converter"
	"github.com/murilo-bracero/raspstore/file-info-service/internal/usecase"
	"github.com/murilo-bracero/raspstore/file-info-service/proto/v1/file-info-service/pb"
)

type UploadHandler interface {
	Upload(w http.ResponseWriter, r *http.Request)
}

type uploadHandler struct {
	uploadUseCase     usecase.UploadFileUseCase
	createFileUseCase usecase.CreateFileUseCase
}

func NewUploadHandler(uploadUseCase usecase.UploadFileUseCase, createFileUseCase usecase.CreateFileUseCase) UploadHandler {
	return &uploadHandler{uploadUseCase: uploadUseCase, createFileUseCase: createFileUseCase}
}

func (h *uploadHandler) Upload(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(rMiddleware.UserIdKey).(string)
	traceId := r.Context().Value(middleware.RequestIDKey).(string)
	r.ParseMultipartForm(32 << 20)

	file, header, err := r.FormFile("file")

	if err != nil {
		log.Printf("[ERROR] - [%s]: Could not open Multipart Form file: %s", traceId, err.Error())
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
		return
	}

	defer file.Close()

	fm := converter.ToFile(&pb.CreateFileMetadataRequest{
		OwnerId:  userId,
		Filename: header.Filename,
		Size:     header.Size,
	})

	if err := h.uploadUseCase.Execute(r.Context(), fm, file); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if err := h.createFileUseCase.Execute(fm); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	u.Created(w, &v1.UploadSuccessResponse{
		FileId:   fm.FileId,
		Filename: fm.Filename,
		OwnerId:  fm.Owner,
	})
}
