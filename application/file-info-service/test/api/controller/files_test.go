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
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"raspstore.github.io/file-manager/api/controller"
	"raspstore.github.io/file-manager/api/dto"
	"raspstore.github.io/file-manager/model"
)

func TestGetAllFilesSuccess(t *testing.T) {
	repo := &filesRepositoryMock{}
	ctr := controller.NewFilesController(repo)

	req, _ := http.NewRequest("GET", "/files", nil)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ctr.ListFiles)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestGetAllFilesPaginatedSuccess(t *testing.T) {
	repo := &filesRepositoryMock{}
	ctr := controller.NewFilesController(repo)

	page := 0
	size := 3

	req, _ := http.NewRequest("GET", fmt.Sprintf("/files?page=%d&size=%d", page, size), nil)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ctr.ListFiles)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestDeleteFileSuccess(t *testing.T) {
	repo := &filesRepositoryMock{}
	ctr := controller.NewFilesController(repo)

	random := uuid.NewString()
	req, _ := http.NewRequest("DELETE", "/files/"+random, nil)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ctr.ListFiles)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)
}

func TestDeleteFileInternalServerError(t *testing.T) {
	repo := &filesRepositoryMock{shouldReturnErr: true}
	ctr := controller.NewFilesController(repo)

	random := uuid.NewString()
	req, _ := http.NewRequest("DELETE", "/files/"+random, nil)
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), middleware.RequestIDKey, "test-trace-id")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ctr.ListFiles)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	var errRes dto.ErrorResponse
	json.Unmarshal(rr.Body.Bytes(), &errRes)

	assert.NotEmpty(t, errRes.Code)
	assert.NotEmpty(t, errRes.Message)
	assert.Equal(t, "test-trace-id", errRes.TraceId)
}

func TestUpdateFileSuccess(t *testing.T) {
	repo := &filesRepositoryMock{}
	ctr := controller.NewFilesController(repo)

	random := uuid.NewString()
	reqBody := []byte(fmt.Sprintf(`{
		"filename": "%s"
	  }`, random))
	req, _ := http.NewRequest("PUT", "/files/"+random, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ctr.ListFiles)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestUpdateFileNotFound(t *testing.T) {
	repo := &filesRepositoryMock{shouldNotFindFile: true}
	ctr := controller.NewFilesController(repo)

	random := uuid.NewString()
	reqBody := []byte(fmt.Sprintf(`{
		"filename": "%s"
	  }`, random))
	req, _ := http.NewRequest("PUT", "/files/"+random, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), middleware.RequestIDKey, "test-trace-id")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ctr.ListFiles)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestUpdateFileInternalServerError(t *testing.T) {
	repo := &filesRepositoryMock{shouldReturnErr: true}
	ctr := controller.NewFilesController(repo)

	random := uuid.NewString()
	reqBody := []byte(fmt.Sprintf(`{
		"filename": "%s"
	  }`, random))
	req, _ := http.NewRequest("PUT", "/files/"+random, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), middleware.RequestIDKey, "test-trace-id")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ctr.ListFiles)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	var errRes dto.ErrorResponse
	json.Unmarshal(rr.Body.Bytes(), &errRes)

	assert.NotEmpty(t, errRes.Code)
	assert.NotEmpty(t, errRes.Message)
	assert.Equal(t, "test-trace-id", errRes.TraceId)
}

type filesRepositoryMock struct {
	shouldReturnErr   bool
	shouldNotFindFile bool
	totalElements     int
}

func (f *filesRepositoryMock) Save(file *model.File) error {
	return errors.New("Not Implemented!")
}

func (f *filesRepositoryMock) FindById(id string) (*model.File, error) {
	if f.shouldReturnErr {
		return nil, mongo.ErrClientDisconnected
	}

	if f.shouldNotFindFile {
		return nil, mongo.ErrNoDocuments
	}

	return createFileMetadata(id), nil
}

func (f *filesRepositoryMock) Delete(id string) error {
	if f.shouldReturnErr {
		return mongo.ErrClientDisconnected
	}

	return nil
}

func (f *filesRepositoryMock) Update(file *model.File) error {
	if f.shouldReturnErr {
		return mongo.ErrClientDisconnected
	}

	return nil
}

func (f *filesRepositoryMock) FindAll(page int, size int) (filesPage *model.FilePage, err error) {
	files := make([]*model.File, 1)
	for i := 0; i < f.totalElements; i++ {
		files = append(files, createFileMetadata(""))
	}

	return &model.FilePage{
		Content: files,
		Count:   size,
	}, nil
}

func createFileMetadata(id string) *model.File {

	if id == "" {
		id = uuid.NewString()
	}

	return &model.File{
		FileId:    id,
		Filename:  id,
		Path:      uuid.NewString() + "/" + id,
		Size:      int64(rand.Int()),
		UpdatedAt: time.Now(),
		CreatedBy: uuid.NewString(),
		UpdatedBy: uuid.NewString(),
	}
}
