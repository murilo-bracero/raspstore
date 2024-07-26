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

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/murilo-bracero/raspstore/file-service/internal/application/facade/mocks"
	"github.com/murilo-bracero/raspstore/file-service/internal/application/repository"
	"github.com/murilo-bracero/raspstore/file-service/internal/domain/entity"
	apiHandler "github.com/murilo-bracero/raspstore/file-service/internal/infra/handler"
	m "github.com/murilo-bracero/raspstore/file-service/internal/infra/middleware"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestGetAllFilesSuccess(t *testing.T) {
	token := jwt.New()
	token.Set("sub", "userId")

	mockCtrl := gomock.NewController(t)

	ff := mocks.NewMockFileFacade(mockCtrl)

	ff.EXPECT().FindAll(gomock.Any(), gomock.Any(), 0, 0, "", false).Return(&entity.FilePage{
		Content: []*entity.File{},
		Count:   0,
	}, nil)

	ctr := apiHandler.NewFilesHandler(ff, nil)

	req, _ := http.NewRequest("GET", "/files", nil)
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), middleware.RequestIDKey, "test-trace-id")
	ctx = context.WithValue(ctx, m.UserClaimsCtxKey, token)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ctr.ListFiles)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestGetAllFilesPaginatedSuccess(t *testing.T) {
	token := jwt.New()
	token.Set("sub", "userId")

	mockCtrl := gomock.NewController(t)

	ff := mocks.NewMockFileFacade(mockCtrl)

	ff.EXPECT().FindAll(gomock.Any(), gomock.Any(), 0, 3, "", false).Return(&entity.FilePage{
		Content: []*entity.File{},
		Count:   0,
	}, nil)

	ctr := apiHandler.NewFilesHandler(ff, nil)

	page := 0
	size := 3

	req, _ := http.NewRequest("GET", fmt.Sprintf("/files?page=%d&size=%d", page, size), nil)
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), middleware.RequestIDKey, "test-trace-id")
	ctx = context.WithValue(ctx, m.UserClaimsCtxKey, token)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ctr.ListFiles)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestGetAllFilesPaginatedInternalServerError(t *testing.T) {
	token := jwt.New()
	token.Set("sub", "userId")

	mockCtrl := gomock.NewController(t)

	ff := mocks.NewMockFileFacade(mockCtrl)

	ff.EXPECT().FindAll(gomock.Any(), gomock.Any(), 0, 3, "", false).Return(nil, errors.New("generic error"))

	ctr := apiHandler.NewFilesHandler(ff, nil)

	page := 0
	size := 3

	req, _ := http.NewRequest("GET", fmt.Sprintf("/files?page=%d&size=%d", page, size), nil)
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), middleware.RequestIDKey, "test-trace-id")
	ctx = context.WithValue(ctx, m.UserClaimsCtxKey, token)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ctr.ListFiles)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestDeleteFileSuccess(t *testing.T) {
	token := jwt.New()
	token.Set("sub", "userId")

	random := uuid.NewString()

	mockCtrl := gomock.NewController(t)

	ff := mocks.NewMockFileFacade(mockCtrl)

	ff.EXPECT().DeleteById(gomock.Any(), gomock.Any(), random).Return(nil)

	ctr := apiHandler.NewFilesHandler(ff, nil)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", random)

	req, _ := http.NewRequest("DELETE", "/files/"+random, nil)
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), middleware.RequestIDKey, "test-trace-id")
	ctx = context.WithValue(ctx, m.UserClaimsCtxKey, token)
	ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ctr.Delete)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)
}

func TestDeleteFileInternalServerError(t *testing.T) {
	token := jwt.New()
	token.Set("sub", "userId")

	random := uuid.NewString()

	mockCtrl := gomock.NewController(t)

	ff := mocks.NewMockFileFacade(mockCtrl)

	ff.EXPECT().DeleteById("test-trace-id", "userId", random).Return(errors.New("generic error"))

	ctr := apiHandler.NewFilesHandler(ff, nil)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", random)

	req, _ := http.NewRequest("DELETE", "/files/"+random, nil)
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), middleware.RequestIDKey, "test-trace-id")
	ctx = context.WithValue(ctx, m.UserClaimsCtxKey, token)
	ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ctr.Delete)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestUpdateFileSuccess(t *testing.T) {
	uc := &updateUseCaseMock{}
	ctr := apiHandler.NewFilesHandler(nil, uc)

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
	ctr := apiHandler.NewFilesHandler(nil, uc)

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
	ctr := apiHandler.NewFilesHandler(nil, uc)

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

type updateUseCaseMock struct {
	shouldThrowError    bool
	shouldThrowNotFound bool
}

func (c *updateUseCaseMock) Execute(ctx context.Context, file *entity.File) (fileMetadata *entity.File, err error) {
	if c.shouldThrowError {
		return nil, errors.New("generic error")
	}

	if c.shouldThrowNotFound {
		return nil, repository.ErrFileDoesNotExists
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
		UpdatedAt: &[]time.Time{time.Now()}[0],
		CreatedBy: uuid.NewString(),
		UpdatedBy: &[]string{uuid.NewString()}[0],
	}
}
