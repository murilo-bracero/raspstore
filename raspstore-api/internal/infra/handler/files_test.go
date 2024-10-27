package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/murilo-bracero/raspstore/file-service/internal/application/facade/mocks"
	"github.com/murilo-bracero/raspstore/file-service/internal/domain/entity"
	"github.com/murilo-bracero/raspstore/file-service/internal/infra/handler"
	"github.com/murilo-bracero/raspstore/file-service/internal/infra/repository"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestGetAllFiles(t *testing.T) {
	token := jwt.New()
	err := token.Set("sub", "userId")
	assert.NoError(t, err)

	mockCtrl := gomock.NewController(t)

	newHandler := func(ff *mocks.MockFileFacade) *handler.Handler {
		ctr := handler.New(nil, ff, nil, nil, nil)
		return ctr
	}

	t.Run("TestGetAllFilesSuccess", func(t *testing.T) {
		ff := mocks.NewMockFileFacade(mockCtrl)

		ff.EXPECT().FindAll(gomock.Any(), gomock.Any(), 0, 0, "", false).Return(&entity.FilePage{
			Content: []*entity.File{},
			Count:   0,
		}, nil)

		ctr := newHandler(ff)

		req, _ := http.NewRequest("GET", "/files", nil)
		req.Header.Set("Content-Type", "application/json")
		ctx := context.WithValue(req.Context(), middleware.RequestIDKey, "test-trace-id")
		ctx = context.WithValue(ctx, handler.UserClaimsCtxKey, token)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(ctr.ListFiles)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("TestGetAllFilesPaginatedSuccess", func(t *testing.T) {
		ff := mocks.NewMockFileFacade(mockCtrl)

		ff.EXPECT().FindAll(gomock.Any(), gomock.Any(), 0, 3, "", false).Return(&entity.FilePage{
			Content: []*entity.File{},
			Count:   0,
		}, nil)

		ctr := newHandler(ff)

		page := 0
		size := 3

		req, _ := http.NewRequest("GET", fmt.Sprintf("/files?page=%d&size=%d", page, size), nil)
		req.Header.Set("Content-Type", "application/json")
		ctx := context.WithValue(req.Context(), middleware.RequestIDKey, "test-trace-id")
		ctx = context.WithValue(ctx, handler.UserClaimsCtxKey, token)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(ctr.ListFiles)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("TestGetAllFilesPaginatedInternalServerError", func(t *testing.T) {
		ff := mocks.NewMockFileFacade(mockCtrl)

		ff.EXPECT().FindAll(gomock.Any(), gomock.Any(), 0, 3, "", false).Return(nil, errors.New("generic error"))

		ctr := newHandler(ff)

		page := 0
		size := 3

		req, _ := http.NewRequest("GET", fmt.Sprintf("/files?page=%d&size=%d", page, size), nil)
		req.Header.Set("Content-Type", "application/json")
		ctx := context.WithValue(req.Context(), middleware.RequestIDKey, "test-trace-id")
		ctx = context.WithValue(ctx, handler.UserClaimsCtxKey, token)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(ctr.ListFiles)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}

func TestDelete(t *testing.T) {
	token := jwt.New()
	err := token.Set("sub", "userId")
	assert.NoError(t, err)

	mockCtrl := gomock.NewController(t)

	newHandler := func(ff *mocks.MockFileFacade) *handler.Handler {
		ctr := handler.New(nil, ff, nil, nil, nil)
		return ctr
	}

	t.Run("TestDeleteFileSuccess", func(t *testing.T) {
		random := uuid.NewString()

		ff := mocks.NewMockFileFacade(mockCtrl)

		ff.EXPECT().DeleteById(gomock.Any(), gomock.Any(), random).Return(nil)

		ctr := newHandler(ff)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", random)

		req, _ := http.NewRequest("DELETE", "/files/"+random, nil)
		req.Header.Set("Content-Type", "application/json")
		ctx := context.WithValue(req.Context(), middleware.RequestIDKey, "test-trace-id")
		ctx = context.WithValue(ctx, handler.UserClaimsCtxKey, token)
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(ctr.Delete)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNoContent, rr.Code)
	})

	t.Run("TestDeleteFileInternalServerError", func(t *testing.T) {
		random := uuid.NewString()

		ff := mocks.NewMockFileFacade(mockCtrl)

		ff.EXPECT().DeleteById("test-trace-id", "userId", random).Return(errors.New("generic error"))

		ctr := newHandler(ff)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", random)

		req, _ := http.NewRequest("DELETE", "/files/"+random, nil)
		req.Header.Set("Content-Type", "application/json")
		ctx := context.WithValue(req.Context(), middleware.RequestIDKey, "test-trace-id")
		ctx = context.WithValue(ctx, handler.UserClaimsCtxKey, token)
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(ctr.Delete)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}

func TestUpdate(t *testing.T) {
	newHandler := func(m *updateUseCaseMock) *handler.Handler {
		ctr := handler.New(m, nil, nil, nil, nil)
		return ctr
	}

	t.Run("TestUpdateFileSuccess", func(t *testing.T) {
		uc := &updateUseCaseMock{}
		ctr := newHandler(uc)

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
		ctx = context.WithValue(ctx, handler.UserClaimsCtxKey, jwt.New())
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
	})

	t.Run("TestUpdateFileNotFound", func(t *testing.T) {
		uc := &updateUseCaseMock{shouldThrowNotFound: true}
		ctr := newHandler(uc)

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
		ctx = context.WithValue(ctx, handler.UserClaimsCtxKey, jwt.New())
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(ctr.Update)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNotFound, rr.Code)
	})

	t.Run("TestUpdateFileInternalServerError", func(t *testing.T) {
		uc := &updateUseCaseMock{shouldThrowError: true}
		ctr := newHandler(uc)

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
		ctx = context.WithValue(ctx, handler.UserClaimsCtxKey, jwt.New())
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(ctr.Update)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}

type updateUseCaseMock struct {
	shouldThrowError    bool
	shouldThrowNotFound bool
}

func (c *updateUseCaseMock) Execute(file *entity.File, userId, traceId string) (fileMetadata *entity.File, err error) {
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
		Size:      13293,
		Owner:     uuid.NewString(),
		Editors:   []string{uuid.NewString(), uuid.NewString(), uuid.NewString()},
		Viewers:   []string{uuid.NewString(), uuid.NewString(), uuid.NewString()},
		CreatedAt: time.Now(),
		UpdatedAt: &[]time.Time{time.Now()}[0],
		CreatedBy: uuid.NewString(),
		UpdatedBy: &[]string{uuid.NewString()}[0],
	}
}
