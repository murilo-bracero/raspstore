package controller_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	rMiddleware "github.com/murilo-bracero/raspstore/commons/pkg/middleware"
	api "github.com/murilo-bracero/raspstore/file-info-service/api/v1"
	"github.com/murilo-bracero/raspstore/file-info-service/internal"
	apiHandler "github.com/murilo-bracero/raspstore/file-info-service/internal/api/handler"
	"github.com/murilo-bracero/raspstore/file-info-service/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestGetAllFilesSuccess(t *testing.T) {
	uc := &listUseCaseMock{}
	ctr := apiHandler.NewFilesHandler(uc, nil, nil)

	req, _ := http.NewRequest("GET", "/files", nil)
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), rMiddleware.UserIdKey, "random-uuid")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ctr.ListFiles)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestGetAllFilesPaginatedSuccess(t *testing.T) {
	uc := &listUseCaseMock{}
	ctr := apiHandler.NewFilesHandler(uc, nil, nil)

	page := 0
	size := 3

	req, _ := http.NewRequest("GET", fmt.Sprintf("/files?page=%d&size=%d", page, size), nil)
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), rMiddleware.UserIdKey, "random-uuid")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ctr.ListFiles)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestGetAllFilesPaginatedInternalServerError(t *testing.T) {
	uc := &listUseCaseMock{shouldThrowError: true}
	ctr := apiHandler.NewFilesHandler(uc, nil, nil)

	page := 0
	size := 3

	req, _ := http.NewRequest("GET", fmt.Sprintf("/files?page=%d&size=%d", page, size), nil)
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), rMiddleware.UserIdKey, "random-uuid")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ctr.ListFiles)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestDeleteFileSuccess(t *testing.T) {
	uc := &deleteUseCaseMock{}
	ctr := apiHandler.NewFilesHandler(nil, nil, uc)

	random := uuid.NewString()
	req, _ := http.NewRequest("DELETE", "/files/"+random, nil)
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), rMiddleware.UserIdKey, "random-uuid")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ctr.Delete)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)
}

func TestDeleteFileInternalServerError(t *testing.T) {
	uc := &deleteUseCaseMock{shouldThrowError: true}
	ctr := apiHandler.NewFilesHandler(nil, nil, uc)

	random := uuid.NewString()
	req, _ := http.NewRequest("DELETE", "/files/"+random, nil)
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), rMiddleware.UserIdKey, "random-uuid")
	ctx = context.WithValue(ctx, middleware.RequestIDKey, "test-trace-id")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ctr.Delete)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	var errRes api.ErrorResponse
	json.Unmarshal(rr.Body.Bytes(), &errRes)

	assert.NotEmpty(t, errRes.Code)
	assert.NotEmpty(t, errRes.Message)
	assert.Equal(t, "test-trace-id", errRes.TraceId)
}

func TestUpdateFileSuccess(t *testing.T) {
	uc := &updateUseCaseMock{}
	ctr := apiHandler.NewFilesHandler(nil, uc, nil)

	random := uuid.NewString()
	reqBody := []byte(fmt.Sprintf(`{
		"filename": "%s"
	  }`, random))
	req, _ := http.NewRequest("PUT", "/files/"+random, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), rMiddleware.UserIdKey, "random-uuid")
	ctx = context.WithValue(ctx, middleware.RequestIDKey, "test-trace-id")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ctr.Update)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var res model.FileMetadataLookup
	err := json.Unmarshal(rr.Body.Bytes(), &res)

	assert.NoError(t, err)

	assert.NotEmpty(t, res.CreatedAt)
	assert.NotEmpty(t, res.UpdatedAt)
	assert.NotEmpty(t, res.CreatedBy)
	assert.NotEmpty(t, res.UpdatedBy)
	assert.NotNil(t, res.Editors)
	assert.NotNil(t, res.Viewers)
	assert.NotEmpty(t, res.FileId)
	assert.NotEmpty(t, res.Filename)
	assert.NotEmpty(t, res.Owner)
	assert.NotNil(t, res.Folder)
	assert.NotEqual(t, 0, res.Size)
}

func TestUpdateFileNotFound(t *testing.T) {
	uc := &updateUseCaseMock{shouldThrowNotFound: true}
	ctr := apiHandler.NewFilesHandler(nil, uc, nil)

	random := uuid.NewString()
	reqBody := []byte(fmt.Sprintf(`{
		"filename": "%s"
	  }`, random))
	req, _ := http.NewRequest("PUT", "/files/"+random, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), middleware.RequestIDKey, "test-trace-id")
	ctx = context.WithValue(ctx, rMiddleware.UserIdKey, "random-uuid")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ctr.Update)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestUpdateFileInternalServerError(t *testing.T) {
	uc := &updateUseCaseMock{shouldThrowError: true}
	ctr := apiHandler.NewFilesHandler(nil, uc, nil)

	random := uuid.NewString()
	reqBody := []byte(fmt.Sprintf(`{
		"filename": "%s"
	  }`, random))
	req, _ := http.NewRequest("PUT", "/files/"+random, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), middleware.RequestIDKey, "test-trace-id")
	ctx = context.WithValue(ctx, rMiddleware.UserIdKey, "random-uuid")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ctr.Update)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	var errRes api.ErrorResponse
	json.Unmarshal(rr.Body.Bytes(), &errRes)

	assert.NotEmpty(t, errRes.Code)
	assert.NotEmpty(t, errRes.Message)
	assert.Equal(t, "test-trace-id", errRes.TraceId)
}

type deleteUseCaseMock struct {
	shouldThrowError bool
}

func (c *deleteUseCaseMock) Execute(ctx context.Context, fileId string) (error_ error) {
	if c.shouldThrowError {
		return errors.New("generic error")
	}

	return
}

type listUseCaseMock struct {
	shouldThrowError bool
}

func (c *listUseCaseMock) Execute(ctx context.Context, page int, size int) (filesPage *model.FilePage, error_ error) {
	if c.shouldThrowError {
		return nil, errors.New("generic error")
	}

	return
}

type updateUseCaseMock struct {
	shouldThrowError    bool
	shouldThrowNotFound bool
}

func (c *updateUseCaseMock) Execute(ctx context.Context, file *model.File) (fileMetadata *model.FileMetadataLookup, error_ error) {
	if c.shouldThrowError {
		return nil, errors.New("generic error")
	}

	if c.shouldThrowNotFound {
		return nil, internal.ErrFileDoesNotExists
	}

	return createFileMetadataLookup(file.FileId), nil
}

func createFileMetadataLookup(id string) *model.FileMetadataLookup {

	if id == "" {
		id = uuid.NewString()
	}

	return &model.FileMetadataLookup{
		FileId:   id,
		Filename: id,
		Folder: model.Folder{
			Name:     "/",
			IsSecret: false,
		},
		Size:      int64(rand.Int()),
		Owner:     model.UserView{UserId: uuid.NewString(), Username: uuid.NewString()},
		Editors:   []model.UserView{{UserId: uuid.NewString(), Username: uuid.NewString()}, {UserId: uuid.NewString(), Username: uuid.NewString()}},
		Viewers:   []model.UserView{{UserId: uuid.NewString(), Username: uuid.NewString()}, {UserId: uuid.NewString(), Username: uuid.NewString()}},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		CreatedBy: model.UserView{UserId: uuid.NewString(), Username: uuid.NewString()},
		UpdatedBy: model.UserView{UserId: uuid.NewString(), Username: uuid.NewString()},
	}
}
