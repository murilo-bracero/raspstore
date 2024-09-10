package handler_test

import (
	"encoding/base64"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/murilo-bracero/raspstore/file-service/internal/application/usecase/mocks"
	"github.com/murilo-bracero/raspstore/file-service/internal/infra/handler"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestAuthenticate(t *testing.T) {
	createReq := func() (req *http.Request) {
		req, err := http.NewRequest("POST", "/file-service/v1/login", nil)
		assert.NoError(t, err)
		return req
	}

	t.Run("happy path", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)

		lpuc := mocks.NewMockLoginPAMUseCase(mockCtrl)

		ctr := handler.New(nil, nil, lpuc, nil, nil, nil)

		creds := base64.StdEncoding.EncodeToString([]byte("username:password"))

		lpuc.EXPECT().Execute("username", "password").Return("token", nil)

		req := createReq()
		req.Header.Set("Authorization", "Basic "+creds)
		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(ctr.Authenticate)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		jar := rr.Result().Cookies()

		assert.NotNil(t, jar, "Cookie Jar is nil")

		ck := jar[0]

		assert.Equal(t, "JWT-TOKEN", ck.Name)
		assert.Equal(t, "token", ck.Value)
		assert.True(t, ck.HttpOnly)
		assert.True(t, ck.Secure)
	})

	t.Run("should return 401 Unauthenticated if usecase returns error", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)

		lpuc := mocks.NewMockLoginPAMUseCase(mockCtrl)

		ctr := handler.New(nil, nil, lpuc, nil, nil, nil)

		creds := base64.StdEncoding.EncodeToString([]byte("username:password"))

		lpuc.EXPECT().Execute("username", "password").Return("", errors.New("generic error"))

		req := createReq()
		req.Header.Set("Authorization", "Basic "+creds)
		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(ctr.Authenticate)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)

		jar := rr.Result().Cookies()

		assert.Empty(t, jar)
	})

	t.Run("should return 401 Unauthenticated if credentials not present", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)

		lpuc := mocks.NewMockLoginPAMUseCase(mockCtrl)

		ctr := handler.New(nil, nil, lpuc, nil, nil, nil)

		req := createReq()
		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(ctr.Authenticate)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)

		jar := rr.Result().Cookies()

		assert.Empty(t, jar)
	})
}
