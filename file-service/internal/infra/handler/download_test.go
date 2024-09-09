package handler_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	chim "github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/murilo-bracero/raspstore/file-service/internal/application/facade/mocks"
	"github.com/murilo-bracero/raspstore/file-service/internal/domain/entity"
	"github.com/murilo-bracero/raspstore/file-service/internal/infra/handler"
	m "github.com/murilo-bracero/raspstore/file-service/internal/infra/middleware"
	"github.com/murilo-bracero/raspstore/file-service/internal/infra/repository"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestDownload(t *testing.T) {
	token := jwt.New()
	err := token.Set("sub", defaultUserId)
	assert.NoError(t, err)

	createReq := func() (req *http.Request) {
		req, _ = http.NewRequest("GET", "/file-service/v1/downloads/4e2bc94b-a6b6-4c44-9512-79b5eb654524", nil)
		ctx := context.WithValue(req.Context(), m.UserClaimsCtxKey, token)
		ctx = context.WithValue(ctx, chim.RequestIDKey, "trace-id")
		return req.WithContext(ctx)
	}

	t.Run("happy path", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)

		ff := mocks.NewMockFileFacade(mockCtrl)
		ffc := mocks.NewMockFileSystemFacade(mockCtrl)

		ff.EXPECT().FindById(gomock.Any(), gomock.Any()).Return(&entity.File{
			FileId:    "4e2bc94b-a6b6-4c44-9512-79b5eb654524",
			Filename:  testFilename,
			Size:      1024,
			UpdatedAt: &[]time.Time{time.Now()}[0],
			CreatedBy: uuid.NewString(),
			UpdatedBy: &[]string{uuid.NewString()}[0],
		}, nil)

		tempFile, err := createTempFile()

		assert.NoError(t, err, "createTempFile")

		fileMock, err := os.Open(tempFile)

		assert.NoError(t, err, "os.Open")

		ffc.EXPECT().Download("trace-id", defaultUserId, "4e2bc94b-a6b6-4c44-9512-79b5eb654524").Return(fileMock, nil)

		ctr := handler.New(nil, nil, nil, ff, ffc, nil)

		req := createReq()

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(ctr.Download)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, "application/octet-stream", rr.Header().Get("Content-Type"))
		assert.Equal(t, fmt.Sprintf("attachment; filename=\"%s\"", testFilename), rr.Header().Get("Content-Disposition"))
	})

	t.Run("should return NOT FOUND when no file are found in database with given id", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)

		ff := mocks.NewMockFileFacade(mockCtrl)
		ffc := mocks.NewMockFileSystemFacade(mockCtrl)

		ff.EXPECT().FindById(gomock.Any(), gomock.Any()).Return(nil, repository.ErrFileDoesNotExists)

		ctr := handler.New(nil, nil, nil, ff, ffc, nil)

		req := createReq()

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(ctr.Download)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNotFound, rr.Code)
	})

	t.Run("should return internal server error when use case returns unhandled error", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)

		ff := mocks.NewMockFileFacade(mockCtrl)
		ffc := mocks.NewMockFileSystemFacade(mockCtrl)

		ff.EXPECT().FindById(gomock.Any(), gomock.Any()).Return(nil, errors.New("generic error"))

		ctr := handler.New(nil, nil, nil, ff, ffc, nil)

		req := createReq()

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(ctr.Download)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}
