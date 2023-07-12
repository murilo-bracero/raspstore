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
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	api "raspstore.github.io/file-manager/api/v1"
	apiHandler "raspstore.github.io/file-manager/internal/api/handler"
	"raspstore.github.io/file-manager/internal/model"
)

func TestGetAllFilesSuccess(t *testing.T) {
	repo := &filesRepositoryMock{}
	ctr := apiHandler.NewFilesHandler(repo)

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
	repo := &filesRepositoryMock{}
	ctr := apiHandler.NewFilesHandler(repo)

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

func TestDeleteFileSuccess(t *testing.T) {
	repo := &filesRepositoryMock{}
	ctr := apiHandler.NewFilesHandler(repo)

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
	repo := &filesRepositoryMock{shouldReturnErr: true}
	ctr := apiHandler.NewFilesHandler(repo)

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
	repo := &filesRepositoryMock{}
	ctr := apiHandler.NewFilesHandler(repo)

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
	assert.NotEmpty(t, res.Path)
	assert.NotEqual(t, 0, res.Size)
}

func TestUpdateFileNotFound(t *testing.T) {
	repo := &filesRepositoryMock{shouldNotFindFile: true}
	ctr := apiHandler.NewFilesHandler(repo)

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
	repo := &filesRepositoryMock{shouldReturnErr: true}
	ctr := apiHandler.NewFilesHandler(repo)

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

type filesRepositoryMock struct {
	shouldReturnErr   bool
	shouldNotFindFile bool
	totalElements     int
}

func (f *filesRepositoryMock) Save(file *model.File) error {
	return errors.New("Not Implemented!")
}

func (f *filesRepositoryMock) FindById(userId string, id string) (*model.File, error) {
	if f.shouldReturnErr {
		return nil, mongo.ErrClientDisconnected
	}

	if f.shouldNotFindFile {
		return nil, mongo.ErrNoDocuments
	}

	return createFileMetadata(id), nil
}

func (f *filesRepositoryMock) FindByIdLookup(userId string, id string) (fileMetadata *model.FileMetadataLookup, err error) {
	if f.shouldReturnErr {
		return nil, mongo.ErrClientDisconnected
	}

	if f.shouldNotFindFile {
		return nil, mongo.ErrNoDocuments
	}

	return createFileMetadataLookup(id), nil
}

func (f *filesRepositoryMock) Delete(userId string, fileId string) error {
	if f.shouldReturnErr {
		return mongo.ErrClientDisconnected
	}

	return nil
}

func (f *filesRepositoryMock) Update(userId string, file *model.File) error {
	if f.shouldReturnErr {
		return mongo.ErrClientDisconnected
	}

	return nil
}

func (f *filesRepositoryMock) FindAll(userId string, page int, size int) (filesPage *model.FilePage, err error) {
	files := make([]*model.FileMetadataLookup, 1)
	for i := 0; i < f.totalElements; i++ {
		files = append(files, createFileMetadataLookup(""))
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

func createFileMetadataLookup(id string) *model.FileMetadataLookup {

	if id == "" {
		id = uuid.NewString()
	}

	return &model.FileMetadataLookup{
		FileId:    id,
		Filename:  id,
		Path:      uuid.NewString() + "/" + id,
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
