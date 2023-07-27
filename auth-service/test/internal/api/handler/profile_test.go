package handler_test

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	v1 "github.com/murilo-bracero/raspstore/auth-service/api/v1"
	"github.com/murilo-bracero/raspstore/auth-service/internal"
	"github.com/murilo-bracero/raspstore/auth-service/internal/api/handler"
	"github.com/murilo-bracero/raspstore/auth-service/internal/api/utils"
	"github.com/murilo-bracero/raspstore/auth-service/internal/model"
	"github.com/murilo-bracero/raspstore/auth-service/internal/token"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestGetProfile(t *testing.T) {
	createJsonRequest := func() (*http.Request, error) {
		req, err := http.NewRequest("GET", "/idp/v1/profile", nil)
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		return req, nil
	}

	t.Run("should return token user profile successfull when token is valid in header", func(t *testing.T) {
		req, err := createJsonRequest()
		assert.NoError(t, err)

		id := "10950f72-29ec-49a8-92bc-53003d7237a3"
		permissions := []string{"admin"}

		req.Header.Set("Authorization", generateToken(id, permissions))

		rr := httptest.NewRecorder()

		ph := handler.NewProfileHandler(&mockGetProfileUseCase{}, nil)
		handler := http.HandlerFunc(ph.GetProfile)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var res v1.UserRepresentation
		err = utils.ParseBody(rr.Body, &res)
		assert.NoError(t, err)
		assert.Equal(t, "c223a9f5-7174-4102-aacc-73f03954dde8", res.UserID)
		assert.Equal(t, "cool_username", res.Username)
		assert.Equal(t, true, res.IsMfaEnabled)
		assert.Equal(t, false, res.IsMfaVerified)
		assert.NotEmpty(t, res.CreatedAt)
		assert.NotEmpty(t, res.UpdatedAt)
	})

	t.Run("should return token user profile successfull when token is valid in token", func(t *testing.T) {
		req, err := createJsonRequest()
		assert.NoError(t, err)

		id := "10950f72-29ec-49a8-92bc-53003d7237a3"
		permissions := []string{"admin"}

		req.AddCookie(createAccessTokenCookie(id, permissions))

		rr := httptest.NewRecorder()

		ph := handler.NewProfileHandler(&mockGetProfileUseCase{}, nil)
		handler := http.HandlerFunc(ph.GetProfile)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var res v1.UserRepresentation
		err = utils.ParseBody(rr.Body, &res)
		assert.NoError(t, err)
		assert.Equal(t, "c223a9f5-7174-4102-aacc-73f03954dde8", res.UserID)
		assert.Equal(t, "cool_username", res.Username)
		assert.Equal(t, true, res.IsMfaEnabled)
		assert.Equal(t, false, res.IsMfaVerified)
		assert.NotEmpty(t, res.CreatedAt)
		assert.NotEmpty(t, res.UpdatedAt)
	})

	t.Run("should return unauthorized when no token", func(t *testing.T) {
		req, err := createJsonRequest()
		assert.NoError(t, err)

		ctr := handler.NewProfileHandler(&mockGetProfileUseCase{}, nil)

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(ctr.GetProfile)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("should return unauthorized when token is malformed in cookie", func(t *testing.T) {
		req, err := createJsonRequest()
		assert.NoError(t, err)

		ctr := handler.NewProfileHandler(&mockGetProfileUseCase{}, nil)

		req.AddCookie(createBadTokenCookie())

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(ctr.GetProfile)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("should return unauthorized when token is malformed in header", func(t *testing.T) {
		req, err := createJsonRequest()
		assert.NoError(t, err)

		ctr := handler.NewProfileHandler(&mockGetProfileUseCase{}, nil)

		req.Header.Set("Authorization", "Bearer badToken")

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(ctr.GetProfile)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("should return internal server error when usecase throws error", func(t *testing.T) {
		req, err := createJsonRequest()
		assert.NoError(t, err)

		ctr := handler.NewProfileHandler(&mockGetProfileUseCase{shouldReturnError: true}, nil)

		id := "10950f72-29ec-49a8-92bc-53003d7237a3"
		permissions := []string{"admin"}

		req.Header.Set("Authorization", generateToken(id, permissions))

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
		ctx := context.WithValue(req.Context(), middleware.RequestIDKey, "test-trace-id")
		req = req.WithContext(ctx)
		return req
	}

	t.Run("shoudl return OK when payload and token are valid", func(t *testing.T) {
		req := createJsonRequest(`{
			"username": "coolusername"
		  }`)

		id := "10950f72-29ec-49a8-92bc-53003d7237a3"
		permissions := []string{"admin"}

		req.Header.Set("Authorization", generateToken(id, permissions))

		rr := httptest.NewRecorder()

		ph := handler.NewProfileHandler(nil, &mockUpdateUserUseCase{shouldReturnError: false})
		handler := http.HandlerFunc(ph.UpdateProfile)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("should return BAD REQUEST when payload is invalid", func(t *testing.T) {
		req := createJsonRequest(`{
		  }`)

		id := "10950f72-29ec-49a8-92bc-53003d7237a3"
		permissions := []string{"admin"}

		req.Header.Set("Authorization", generateToken(id, permissions))

		rr := httptest.NewRecorder()

		ph := handler.NewProfileHandler(nil, &mockUpdateUserUseCase{shouldReturnError: false})
		handler := http.HandlerFunc(ph.UpdateProfile)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("should return UNAUTHORIZED when token is invalid in header", func(t *testing.T) {
		req := createJsonRequest(`{
			"username": "coolusername"
		  }`)

		req.Header.Set("Authorization", "Bad token")

		rr := httptest.NewRecorder()

		ph := handler.NewProfileHandler(nil, &mockUpdateUserUseCase{shouldReturnError: false})
		handler := http.HandlerFunc(ph.UpdateProfile)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("should return UNAUTHORIZED when token is invalid in cookie", func(t *testing.T) {
		req := createJsonRequest(`{
			"username": "coolusername"
		  }`)

		req.AddCookie(createBadTokenCookie())

		rr := httptest.NewRecorder()

		ph := handler.NewProfileHandler(nil, &mockUpdateUserUseCase{shouldReturnError: false})
		handler := http.HandlerFunc(ph.UpdateProfile)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("should return CONFLICT when usecase return ErrConflict", func(t *testing.T) {
		req := createJsonRequest(`{
			"username": "coolusername"
		  }`)

		id := "10950f72-29ec-49a8-92bc-53003d7237a3"
		permissions := []string{"admin"}

		req.AddCookie(createAccessTokenCookie(id, permissions))

		rr := httptest.NewRecorder()

		ph := handler.NewProfileHandler(nil, &mockUpdateUserUseCase{shouldReturnConflictError: true})
		handler := http.HandlerFunc(ph.UpdateProfile)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusConflict, rr.Code)
	})

	t.Run("should return INTERNAL SERVER ERROR when usecase return unexpected error", func(t *testing.T) {
		req := createJsonRequest(`{
			"username": "coolusername"
		  }`)

		id := "10950f72-29ec-49a8-92bc-53003d7237a3"
		permissions := []string{"admin"}

		req.AddCookie(createAccessTokenCookie(id, permissions))

		rr := httptest.NewRecorder()

		ph := handler.NewProfileHandler(nil, &mockUpdateUserUseCase{shouldReturnError: true})
		handler := http.HandlerFunc(ph.UpdateProfile)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}

func generateToken(id string, permissions []string) string {
	token, err := token.Generate(&model.User{UserId: id, Permissions: []string{"admin"}})

	if err != nil {
		panic(err)
	}

	return "Bearer " + token
}

func createAccessTokenCookie(id string, permissions []string) *http.Cookie {
	return &http.Cookie{
		Name:     "access_token",
		Value:    generateToken(id, permissions),
		Expires:  time.Now().Add(1 * time.Hour),
		Secure:   true,
		HttpOnly: true,
	}
}

func createBadTokenCookie() *http.Cookie {
	return &http.Cookie{
		Name:     "access_token",
		Value:    "Bearer badToken",
		Expires:  time.Now().Add(1 * time.Hour),
		Secure:   true,
		HttpOnly: true,
	}
}

type mockGetProfileUseCase struct {
	shouldReturnError bool
}

func (u *mockGetProfileUseCase) Execute(ctx context.Context, userId string) (user *model.User, error_ error) {
	if u.shouldReturnError {
		return nil, errors.New("generic error")
	}

	return &model.User{
		Id:            primitive.NewObjectID(),
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
	}, nil
}

type mockUpdateUserUseCase struct {
	shouldReturnError         bool
	shouldReturnConflictError bool
}

func (u *mockUpdateUserUseCase) Execute(ctx context.Context, user *model.User) error {
	if u.shouldReturnError {
		return errors.New("generic error")
	}

	if u.shouldReturnConflictError {
		return internal.ErrConflict
	}

	return nil
}
