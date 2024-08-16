package handler_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/murilo-bracero/raspstore/file-service/internal/domain/entity"
	"github.com/murilo-bracero/raspstore/file-service/internal/infra/config"
	"github.com/murilo-bracero/raspstore/file-service/internal/infra/handler"
	m "github.com/murilo-bracero/raspstore/file-service/internal/infra/middleware"
	"github.com/stretchr/testify/assert"
)

const testFilename = "test.txt"
const defaultUserId = "e9e28c79-a5e8-4545-bd32-e536e690bd4a"

func TestUpload(t *testing.T) {
	config := &config.Config{Storage: struct {
		Path  string
		Limit string
	}{Path: "./"}}

	token := jwt.New()
	token.Set("sub", defaultUserId)

	createReq := func(body *bytes.Buffer) (req *http.Request) {
		req, err := http.NewRequest("POST", "/file-service/v1/uploads", body)
		assert.NoError(t, err)
		ctx := context.WithValue(req.Context(), m.UserClaimsCtxKey, token)
		ctx = context.WithValue(ctx, chiMiddleware.RequestIDKey, defaultUserId)
		req = req.WithContext(ctx)
		return req
	}

	t.Run("happy path", func(t *testing.T) {
		uploadUseCase := &uploadFileUseCaseMock{}
		cFileUseCase := &createUseCaseMock{}
		ctr := handler.NewUploadHandler(config, uploadUseCase, cFileUseCase)

		tempFile, err := createTempFile()
		assert.NoError(t, err)

		file, err := os.Open(tempFile)
		assert.NoError(t, err)
		defer file.Close()

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile("file", filepath.Base(tempFile))
		assert.NoError(t, err)

		_, err = io.Copy(part, file)
		assert.NoError(t, err)

		err = writer.Close()
		assert.NoError(t, err)

		req := createReq(body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(ctr.Upload)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusCreated, rr.Code)
	})

	t.Run("should return bad request when form without file", func(t *testing.T) {
		uploadUseCase := &uploadFileUseCaseMock{}
		cFileUseCase := &createUseCaseMock{}
		ctr := handler.NewUploadHandler(config, uploadUseCase, cFileUseCase)

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		err := writer.Close()
		assert.NoError(t, err)

		req := createReq(body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(ctr.Upload)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
	})

	t.Run("should return bad request when form without file", func(t *testing.T) {
		uploadUseCase := &uploadFileUseCaseMock{}
		cFileUseCase := &createUseCaseMock{}
		ctr := handler.NewUploadHandler(config, uploadUseCase, cFileUseCase)

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		err := writer.Close()
		assert.NoError(t, err)

		req := createReq(body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(ctr.Upload)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
	})

	t.Run("should return internal server error when upload use case returns error", func(t *testing.T) {
		uploadUseCase := &uploadFileUseCaseMock{shouldReturnError: true}
		cFileUseCase := &createUseCaseMock{}
		ctr := handler.NewUploadHandler(config, uploadUseCase, cFileUseCase)

		tempFile, err := createTempFile()
		assert.NoError(t, err)

		file, err := os.Open(tempFile)
		assert.NoError(t, err)
		defer file.Close()

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile("file", filepath.Base(tempFile))
		assert.NoError(t, err)

		_, err = io.Copy(part, file)
		assert.NoError(t, err)

		err = writer.Close()
		assert.NoError(t, err)

		req := createReq(body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(ctr.Upload)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})

	t.Run("should return internal server error when create use case returns error", func(t *testing.T) {
		uploadUseCase := &uploadFileUseCaseMock{}
		cFileUseCase := &createUseCaseMock{shouldReturnErr: true}
		ctr := handler.NewUploadHandler(config, uploadUseCase, cFileUseCase)

		tempFile, err := createTempFile()
		assert.NoError(t, err)

		file, err := os.Open(tempFile)
		assert.NoError(t, err)
		defer file.Close()

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile("file", filepath.Base(tempFile))
		assert.NoError(t, err)

		_, err = io.Copy(part, file)
		assert.NoError(t, err)

		err = writer.Close()
		assert.NoError(t, err)

		req := createReq(body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(ctr.Upload)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}

func createTempFile() (string, error) {
	tempDir := os.TempDir()
	tempFile := filepath.Join(tempDir, testFilename)
	err := os.WriteFile(tempFile, []byte("test content"), 0666)
	return tempFile, err
}

type uploadFileUseCaseMock struct {
	shouldReturnError bool
}

func (u *uploadFileUseCaseMock) Execute(ctx context.Context, file *entity.File, src io.Reader) (err error) {
	if u.shouldReturnError {
		return errors.New("generic error")
	}

	return nil
}

type createUseCaseMock struct {
	shouldReturnErr bool
}

func (f *createUseCaseMock) Execute(file *entity.File) (err error) {
	if f.shouldReturnErr {
		return errors.New("generic error")
	}

	return nil
}
