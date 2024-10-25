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
	"github.com/murilo-bracero/raspstore/file-service/internal/application/facade/mocks"
	"github.com/murilo-bracero/raspstore/file-service/internal/infra/config"
	"github.com/murilo-bracero/raspstore/file-service/internal/infra/handler"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

const testFilename = "test.txt"
const defaultUserId = "e9e28c79-a5e8-4545-bd32-e536e690bd4a"

func TestUpload(t *testing.T) {
	config := &config.Config{Storage: config.StorageConfig{Path: os.TempDir()}}

	token := jwt.New()
	err := token.Set("sub", defaultUserId)
	assert.NoError(t, err)

	createReq := func(body *bytes.Buffer) (req *http.Request) {
		req, err := http.NewRequest("POST", "/file-service/v1/uploads", body)
		assert.NoError(t, err)
		ctx := context.WithValue(req.Context(), handler.UserClaimsCtxKey, token)
		ctx = context.WithValue(ctx, chiMiddleware.RequestIDKey, defaultUserId)
		req = req.WithContext(ctx)
		return req
	}

	t.Run("happy path", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)

		ffc := mocks.NewMockFileSystemFacade(mockCtrl)

		ff := mocks.NewMockFileFacade(mockCtrl)

		ctr := handler.New(nil, ff, ffc, config)

		ff.EXPECT().Save(gomock.Any()).Return(nil)
		ffc.EXPECT().Upload(defaultUserId, gomock.Any(), gomock.Any())

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
		mockCtrl := gomock.NewController(t)

		ffc := mocks.NewMockFileSystemFacade(mockCtrl)

		ctr := handler.New(nil, nil, ffc, config)

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
		mockCtrl := gomock.NewController(t)

		ffc := mocks.NewMockFileSystemFacade(mockCtrl)

		ctr := handler.New(nil, nil, ffc, config)

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
		mockCtrl := gomock.NewController(t)

		ffc := mocks.NewMockFileSystemFacade(mockCtrl)

		ctr := handler.New(nil, nil, ffc, config)

		ffc.EXPECT().Upload(defaultUserId, gomock.Any(), gomock.Any()).Return(errors.New("generic error"))

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

	t.Run("should return internal server error when save file returns error", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)

		ffc := mocks.NewMockFileSystemFacade(mockCtrl)

		ff := mocks.NewMockFileFacade(mockCtrl)

		ctr := handler.New(nil, ff, ffc, config)

		ff.EXPECT().Save(gomock.Any()).Return(errors.New("generic error"))
		ffc.EXPECT().Upload(defaultUserId, gomock.Any(), gomock.Any())

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
	err := os.WriteFile(tempFile, []byte("test content"), 0600)
	return tempFile, err
}
