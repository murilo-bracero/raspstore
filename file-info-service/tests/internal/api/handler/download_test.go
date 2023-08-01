package handler_test

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	rMiddleware "github.com/murilo-bracero/raspstore/commons/pkg/middleware"
	"github.com/murilo-bracero/raspstore/file-service/internal"
	"github.com/murilo-bracero/raspstore/file-service/internal/api/handler"
	"github.com/murilo-bracero/raspstore/file-service/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestDownload(t *testing.T) {
	if err := godotenv.Load("../../test.env"); err != nil {
		log.Println("Could not load .env file. Using system variables instead")
	}

	createReq := func() (req *http.Request) {
		req, _ = http.NewRequest("GET", "/file-service/v1/downloads/4e2bc94b-a6b6-4c44-9512-79b5eb654524", nil)
		ctx := context.WithValue(req.Context(), rMiddleware.UserIdKey, defaultUserId)
		return req.WithContext(ctx)
	}

	t.Run("happy path", func(t *testing.T) {
		downloadUseCase := &downloadFileUseCaseMock{}
		getFileUseCase := &getFileUseCaseMock{}

		ctr := handler.NewDownloadHandler(downloadUseCase, getFileUseCase)

		req := createReq()

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(ctr.Download)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, "application/octet-stream", rr.Header().Get("Content-Type"))
		assert.Equal(t, fmt.Sprintf("attachment; filename=\"%s\"", testFilename), rr.Header().Get("Content-Disposition"))
	})

	t.Run("should return NOT FOUND when no file are found in database with given id", func(t *testing.T) {
		downloadUseCase := &downloadFileUseCaseMock{}
		getFileUseCase := &getFileUseCaseMock{shouldReturnNotFound: true}

		ctr := handler.NewDownloadHandler(downloadUseCase, getFileUseCase)

		req := createReq()

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(ctr.Download)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNotFound, rr.Code)
	})

	t.Run("should return internal server error when use case returns unhandled error", func(t *testing.T) {
		downloadUseCase := &downloadFileUseCaseMock{}
		getFileUseCase := &getFileUseCaseMock{shouldReturnErr: true}

		ctr := handler.NewDownloadHandler(downloadUseCase, getFileUseCase)

		req := createReq()

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(ctr.Download)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}

type downloadFileUseCaseMock struct {
	shouldReturnErr bool
}

func (d *downloadFileUseCaseMock) Execute(ctx context.Context, fileId string) (file *os.File, error_ error) {
	if d.shouldReturnErr {
		return nil, errors.New("generic error")
	}

	tempFile, _ := createTempFile()

	return os.Open(tempFile)
}

type getFileUseCaseMock struct {
	shouldReturnErr      bool
	shouldReturnNotFound bool
}

func (f *getFileUseCaseMock) Execute(userId string, fileId string) (file *model.File, error_ error) {
	if f.shouldReturnErr {
		return nil, errors.New("generic error")
	}

	if f.shouldReturnNotFound {
		return nil, internal.ErrFileDoesNotExists
	}

	return &model.File{
		FileId:    fileId,
		Filename:  testFilename,
		Size:      int64(rand.Int()),
		UpdatedAt: time.Now(),
		CreatedBy: uuid.NewString(),
		UpdatedBy: uuid.NewString(),
	}, nil
}
