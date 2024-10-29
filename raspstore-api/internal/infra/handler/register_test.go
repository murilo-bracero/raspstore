package handler_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/murilo-bracero/raspstore/file-service/internal/application/facade/mocks"
	"github.com/murilo-bracero/raspstore/file-service/internal/infra/config"
	"github.com/murilo-bracero/raspstore/file-service/internal/infra/handler"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestRegister(t *testing.T) {
	createReq := func(body io.Reader) (req *http.Request) {
		req, _ = http.NewRequest("GET", "/file-service/v1/register", body)
		ctx := context.WithValue(req.Context(), middleware.RequestIDKey, "trace-id")
		return req.WithContext(ctx)
	}

	mockConfig := &config.Config{Auth: config.AuthConfig{EnableUserRegister: true}}

	newHandler := func(uf *mocks.MockUserFacade) *handler.Handler {
		ctr := handler.New(nil, nil, uf, nil, mockConfig)
		return ctr
	}

	t.Run("should return created if success", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)

		uf := mocks.NewMockUserFacade(mockCtrl)
		uf.EXPECT().Save(gomock.Any()).Return(nil)

		ctr := newHandler(uf)

		body := `{
			"username": "UsrNm",
			"password": "PssWd1"
		}`

		req := createReq(strings.NewReader(body))

		rr := httptest.NewRecorder()

		endpoint := http.HandlerFunc(ctr.Register)
		endpoint.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusCreated, rr.Code)
	})

	t.Run("should return created with name if success", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)

		uf := mocks.NewMockUserFacade(mockCtrl)
		uf.EXPECT().Save(gomock.Any()).Return(nil)

		ctr := newHandler(uf)

		body := `{
			"username": "UsrNm",
			"password": "PssWd1",
			"name": "My Name"
		}`

		req := createReq(strings.NewReader(body))

		rr := httptest.NewRecorder()

		endpoint := http.HandlerFunc(ctr.Register)
		endpoint.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusCreated, rr.Code)
	})

	t.Run("should return bad request if payload incomplete", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)

		uf := mocks.NewMockUserFacade(mockCtrl)
		uf.EXPECT().Save(gomock.Any()).Times(0)

		ctr := newHandler(uf)

		body := `{
			"password": "PssWd1",
			"name": "My Name"
		}`

		req := createReq(strings.NewReader(body))

		rr := httptest.NewRecorder()

		endpoint := http.HandlerFunc(ctr.Register)
		endpoint.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("should return unprocessable entity if payload malformed", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)

		uf := mocks.NewMockUserFacade(mockCtrl)
		uf.EXPECT().Save(gomock.Any()).Times(0)

		ctr := newHandler(uf)

		body := `
			"password": "PssWd1",
			"name": "My Name"
		}`

		req := createReq(strings.NewReader(body))

		rr := httptest.NewRecorder()

		endpoint := http.HandlerFunc(ctr.Register)
		endpoint.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
	})

	t.Run("should return internal server error if facade returns error", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)

		uf := mocks.NewMockUserFacade(mockCtrl)
		uf.EXPECT().Save(gomock.Any()).Return(errors.New("generic error"))

		ctr := newHandler(uf)

		body := `{
			"username": "UsrNm",
			"password": "PssWd1",
			"name": "My Name"
		}`

		req := createReq(strings.NewReader(body))

		rr := httptest.NewRecorder()

		endpoint := http.HandlerFunc(ctr.Register)
		endpoint.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})

	t.Run("should return unprocessable entity if user registration config is disabled", func(t *testing.T) {
		mockConfig.Auth.EnableUserRegister = false

		mockCtrl := gomock.NewController(t)

		uf := mocks.NewMockUserFacade(mockCtrl)
		uf.EXPECT().Save(gomock.Any()).Times(0)

		ctr := newHandler(uf)

		body := `{
			"username": "UsrNm",
			"password": "PssWd1",
			"name": "My Name"
		}`

		req := createReq(strings.NewReader(body))

		rr := httptest.NewRecorder()

		endpoint := http.HandlerFunc(ctr.Register)
		endpoint.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
	})
}
