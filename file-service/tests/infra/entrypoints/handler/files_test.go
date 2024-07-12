package handler_test

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
	"github.com/murilo-bracero/raspstore/file-service/internal/domain/entity"
	e "github.com/murilo-bracero/raspstore/file-service/internal/domain/exceptions"
	apiHandler "github.com/murilo-bracero/raspstore/file-service/internal/infra/entrypoints/handler"
	"github.com/stretchr/testify/assert"
)

func TestGetAllFilesSuccess(t *testing.T) {
	uc := &listUseCaseMock{}
	ctr := apiHandler.NewFilesHandler(uc, nil, nil)

	req, _ := http.NewRequest("GET", "/files", nil)
	req.Header.Set("Content-Type", "application/json")

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
	ctx := context.WithValue(req.Context(), middleware.RequestIDKey, "test-trace-id")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ctr.Delete)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestUpdateFileSuccess(t *testing.T) {
	uc := &updateUseCaseMock{}
	ctr := apiHandler.NewFilesHandler(nil, uc, nil)

	random := uuid.NewString()
	reqBody := []byte(`{
		"filename": "now_its_secret.docx",
		  "secret": true, 
		  "viewers": ["c74d7720-0026-4466-b59f-d1b4a7f6886f"],
		  "editors": []
	  }`)
	req, _ := http.NewRequest("PUT", "/files/"+random, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), middleware.RequestIDKey, "test-trace-id")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ctr.Update)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var res entity.File
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
	assert.NotEqual(t, 0, res.Size)
}

func TestUpdateFileNotFound(t *testing.T) {
	uc := &updateUseCaseMock{shouldThrowNotFound: true}
	ctr := apiHandler.NewFilesHandler(nil, uc, nil)

	random := uuid.NewString()
	reqBody := []byte(`{
		"filename": "now_its_secret.docx",
		  "secret": true, 
		  "viewers": ["c74d7720-0026-4466-b59f-d1b4a7f6886f"],
		  "editors": []
	  }`)
	req, _ := http.NewRequest("PUT", "/files/"+random, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), middleware.RequestIDKey, "test-trace-id")
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
	reqBody := []byte(`{
			"filename": "now_its_secret.docx",
			"secret": true, 
			"viewers": ["c74d7720-0026-4466-b59f-d1b4a7f6886f"],
			"editors": []
	  }`)
	req, _ := http.NewRequest("PUT", "/files/"+random, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), middleware.RequestIDKey, "test-trace-id")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ctr.Update)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
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

func (c *listUseCaseMock) Execute(ctx context.Context, page int, size int, filename string, secret bool) (filesPage *entity.FilePage, error_ error) {
	if c.shouldThrowError {
		return nil, errors.New("generic error")
	}

	return &entity.FilePage{
		Content: []*entity.File{},
		Count:   0,
	}, nil
}

type updateUseCaseMock struct {
	shouldThrowError    bool
	shouldThrowNotFound bool
}

func (c *updateUseCaseMock) Execute(ctx context.Context, file *entity.File) (fileMetadata *entity.File, error_ error) {
	if c.shouldThrowError {
		return nil, errors.New("generic error")
	}

	if c.shouldThrowNotFound {
		return nil, e.ErrFileDoesNotExists
	}

	return createFileMetadataLookup(file.FileId), nil
}

func createFileMetadataLookup(id string) *entity.File {

	if id == "" {
		id = uuid.NewString()
	}

	return &entity.File{
		FileId:    id,
		Filename:  id,
		Size:      int64(rand.Int()),
		Owner:     uuid.NewString(),
		Editors:   []string{uuid.NewString(), uuid.NewString(), uuid.NewString()},
		Viewers:   []string{uuid.NewString(), uuid.NewString(), uuid.NewString()},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		CreatedBy: uuid.NewString(),
		UpdatedBy: uuid.NewString(),
	}
}