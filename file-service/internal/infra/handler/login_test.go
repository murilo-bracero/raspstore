package handler_test

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/murilo-bracero/raspstore/file-service/internal/application/usecase/mocks"
	"github.com/murilo-bracero/raspstore/file-service/internal/domain/model"
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

		var res model.LoginResponse

		err := json.Unmarshal(rr.Body.Bytes(), &res)

		assert.NoError(t, err)

		assert.Equal(t, "token", res.AccessToken)
		assert.NotEmpty(t, res.ExpiresIn)
		assert.Equal(t, "Bearer", res.Prefix)
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
	})
}
