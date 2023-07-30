package handler_test

import (
	"bytes"
	"context"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	cm "github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/murilo-bracero/raspstore/auth-service/internal"
	"github.com/murilo-bracero/raspstore/auth-service/internal/api/handler"
	"github.com/murilo-bracero/raspstore/auth-service/internal/api/middleware"
	"github.com/murilo-bracero/raspstore/auth-service/internal/infra"
	"github.com/murilo-bracero/raspstore/auth-service/internal/model"
	"github.com/stretchr/testify/assert"
)

var config *infra.Config

func init() {
	err := godotenv.Load("../../../.env.test")

	if err != nil {
		log.Panicln(err.Error())
	}

	config = infra.NewConfig()
}

func TestGetProfile(t *testing.T) {
	createJsonRequest := func() (*http.Request, error) {
		req, err := http.NewRequest("GET", "/idp/v1/profile", nil)
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		ctx := context.WithValue(req.Context(), middleware.UserClaimsCtxKey, &model.UserClaims{})
		req = req.WithContext(ctx)
		return req, nil
	}

	t.Run("happy path", func(t *testing.T) {
		req, err := createJsonRequest()
		assert.NoError(t, err)

		rr := httptest.NewRecorder()

		ph := handler.NewProfileHandler(&mockGetProfileUseCase{}, nil, nil)
		handler := http.HandlerFunc(ph.GetProfile)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("should return FORBIDDEN when user account is not active", func(t *testing.T) {
		req, err := createJsonRequest()
		assert.NoError(t, err)

		ctr := handler.NewProfileHandler(&mockGetProfileUseCase{accountInactive: true}, nil, nil)

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(ctr.GetProfile)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusForbidden, rr.Code)
	})

	t.Run("should return internal server error when usecase throws error", func(t *testing.T) {
		req, err := createJsonRequest()
		assert.NoError(t, err)

		ctr := handler.NewProfileHandler(&mockGetProfileUseCase{shouldReturnError: true}, nil, nil)

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(ctr.GetProfile)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}

func TestUpdateProfile(t *testing.T) {
	createJsonRequest := func(body string) *http.Request {
		reqBody := []byte(body)
		req, err := http.NewRequest("PUT", "/idp/v1/profile", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		ctx := context.WithValue(req.Context(), cm.RequestIDKey, "test-trace-id")
		ctx = context.WithValue(ctx, middleware.UserClaimsCtxKey, &model.UserClaims{})
		req = req.WithContext(ctx)
		return req
	}

	t.Run("happy path - return OK when payload is valid", func(t *testing.T) {
		req := createJsonRequest(`{
			"username": "coolusername"
		  }`)

		rr := httptest.NewRecorder()

		ph := handler.NewProfileHandler(&mockGetProfileUseCase{}, &mockUpdateUserUseCase{}, nil)
		handler := http.HandlerFunc(ph.UpdateProfile)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("should return BAD REQUEST when payload is invalid", func(t *testing.T) {
		req := createJsonRequest(`{
		  }`)

		rr := httptest.NewRecorder()

		ph := handler.NewProfileHandler(&mockGetProfileUseCase{}, &mockUpdateUserUseCase{}, nil)
		handler := http.HandlerFunc(ph.UpdateProfile)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("should return FORBIDDEN when user is inactive", func(t *testing.T) {
		req := createJsonRequest(`{
			"username": "coolusername"
		  }`)

		rr := httptest.NewRecorder()

		ph := handler.NewProfileHandler(&mockGetProfileUseCase{accountInactive: true}, &mockUpdateUserUseCase{}, nil)
		handler := http.HandlerFunc(ph.UpdateProfile)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusForbidden, rr.Code)
	})

	t.Run("should return CONFLICT when usecase return ErrConflict", func(t *testing.T) {
		req := createJsonRequest(`{
			"username": "coolusername"
		  }`)

		rr := httptest.NewRecorder()

		ph := handler.NewProfileHandler(&mockGetProfileUseCase{}, &mockUpdateProfileUseCase{shouldReturnConflictError: true}, nil)
		handler := http.HandlerFunc(ph.UpdateProfile)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusConflict, rr.Code)
	})

	t.Run("should return INTERNAL SERVER ERROR when usecase return unexpected error", func(t *testing.T) {
		req := createJsonRequest(`{
			"username": "coolusername"
		  }`)

		rr := httptest.NewRecorder()

		ph := handler.NewProfileHandler(&mockGetProfileUseCase{}, &mockUpdateProfileUseCase{shouldReturnError: true}, nil)
		handler := http.HandlerFunc(ph.UpdateProfile)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}

func TestDeleteProfile(t *testing.T) {
	createJsonRequest := func() *http.Request {
		req, err := http.NewRequest("DELETE", "/idp/v1/profile", nil)
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		ctx := context.WithValue(req.Context(), cm.RequestIDKey, "test-trace-id")
		ctx = context.WithValue(ctx, middleware.UserClaimsCtxKey, &model.UserClaims{})
		req = req.WithContext(ctx)
		return req
	}

	t.Run("should return OK when token is valid", func(t *testing.T) {
		req := createJsonRequest()

		rr := httptest.NewRecorder()

		ph := handler.NewProfileHandler(&mockGetProfileUseCase{}, nil, &mockDeleteUserUseCase{})
		handler := http.HandlerFunc(ph.DeleteProfile)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNoContent, rr.Code)
	})

	t.Run("should return INTERNAL SERVER ERROR when usecase returns an error", func(t *testing.T) {
		req := createJsonRequest()

		rr := httptest.NewRecorder()

		ph := handler.NewProfileHandler(&mockGetProfileUseCase{}, nil, &mockDeleteUserUseCase{shouldReturnError: true})
		handler := http.HandlerFunc(ph.DeleteProfile)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}

type mockGetProfileUseCase struct {
	shouldReturnError bool
	accountInactive   bool
}

func (u *mockGetProfileUseCase) Execute(ctx context.Context, userId string) (user *model.User, error_ error) {
	if u.shouldReturnError {
		return nil, errors.New("generic error")
	}

	usr := &model.User{
		UserId:        "c223a9f5-7174-4102-aacc-73f03954dde8",
		Username:      "cool_username",
		IsEnabled:     true,
		Password:      "hashed_password",
		Secret:        "user_secret",
		Permissions:   []string{"read", "write"},
		RefreshToken:  "refresh_token",
		IsMfaEnabled:  true,
		IsMfaVerified: false,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if u.accountInactive {
		usr.IsEnabled = false
	}

	return usr, nil
}

type mockUpdateProfileUseCase struct {
	shouldReturnError                bool
	shouldReturnConflictError        bool
	shouldReturnInactiveAccountError bool
}

func (u *mockUpdateProfileUseCase) Execute(ctx context.Context, user *model.User) error {
	if u.shouldReturnError {
		return errors.New("generic error")
	}

	if u.shouldReturnConflictError {
		return internal.ErrConflict
	}

	if u.shouldReturnInactiveAccountError {
		return internal.ErrInactiveAccount
	}

	return nil
}

type mockDeleteUserUseCase struct {
	shouldReturnError bool
}

func (u *mockDeleteUserUseCase) Execute(ctx context.Context, userId string) error {
	if u.shouldReturnError {
		return errors.New("generic error")
	}

	return nil
}
