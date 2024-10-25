package handler_test

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/murilo-bracero/raspstore/file-service/internal/domain/model"
	"github.com/murilo-bracero/raspstore/file-service/internal/infra/config"
	"github.com/murilo-bracero/raspstore/file-service/internal/infra/handler"
	"github.com/stretchr/testify/assert"
)

func TestAuthenticate(t *testing.T) {
	createReq := func() (req *http.Request) {
		req, err := http.NewRequest("POST", "/file-service/v1/login", nil)
		ctx := context.WithValue(req.Context(), middleware.RequestIDKey, "test-trace-id")

		req = req.WithContext(ctx)
		assert.NoError(t, err)
		return req
	}

	t.Run("happy path", func(t *testing.T) {
		ctr := handler.New(nil, nil, nil, nil)

		mockLoginFunc := func(_ *config.Config, username, password string) (string, error) {
			assert.Equal(t, "username", username)
			assert.Equal(t, "password", password)
			return "token", nil
		}

		ctr.LoginFunc = mockLoginFunc

		creds := base64.StdEncoding.EncodeToString([]byte("username:password"))

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

		assert.Equal(t, "test-trace-id", middleware.GetReqID(req.Context()))
	})

	t.Run("should return 401 Unauthenticated if usecase returns error", func(t *testing.T) {
		ctr := handler.New(nil, nil, nil, nil)

		mockLoginFunc := func(_ *config.Config, username, password string) (string, error) {
			return "", errors.New("generic error")
		}

		ctr.LoginFunc = mockLoginFunc

		creds := base64.StdEncoding.EncodeToString([]byte("username:password"))

		req := createReq()
		req.Header.Set("Authorization", "Basic "+creds)
		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(ctr.Authenticate)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("should return 401 Unauthenticated if credentials not present", func(t *testing.T) {
		ctr := handler.New(nil, nil, nil, nil)

		req := createReq()
		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(ctr.Authenticate)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})
}
