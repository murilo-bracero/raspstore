package handler_test

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/murilo-bracero/raspstore/file-service/internal/application/repository"
	"github.com/murilo-bracero/raspstore/file-service/internal/domain/entity"
	"github.com/murilo-bracero/raspstore/file-service/internal/infra/handler"
	m "github.com/murilo-bracero/raspstore/file-service/internal/infra/middleware"
	"github.com/stretchr/testify/assert"
)

func TestDownload(t *testing.T) {
	token := jwt.New()
	token.Set("sub", defaultUserId)

	createReq := func() (req *http.Request) {
		req, _ = http.NewRequest("GET", "/file-service/v1/downloads/4e2bc94b-a6b6-4c44-9512-79b5eb654524", nil)
		ctx := context.WithValue(req.Context(), m.UserClaimsCtxKey, token)
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

func (d *downloadFileUseCaseMock) Execute(ctx context.Context, fileId string) (file *os.File, err error) {
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

func (f *getFileUseCaseMock) Execute(userId string, fileId string) (file *entity.File, err error) {
	if f.shouldReturnErr {
		return nil, errors.New("generic error")
	}

	if f.shouldReturnNotFound {
		return nil, repository.ErrFileDoesNotExists
	}

	return &entity.File{
		FileId:    fileId,
		Filename:  testFilename,
		Size:      int64(rand.Int()),
		UpdatedAt: time.Now(),
		CreatedBy: uuid.NewString(),
		UpdatedBy: uuid.NewString(),
	}, nil
}
