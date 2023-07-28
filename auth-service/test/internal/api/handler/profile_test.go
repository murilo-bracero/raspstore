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
	"github.com/murilo-bracero/raspstore/auth-service/internal/model"
	"github.com/murilo-bracero/raspstore/auth-service/internal/token"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func init() {
	err := godotenv.Load("../../.env.test")

	if err != nil {
		log.Panicln(err.Error())
	}
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

		id := "10950f72-29ec-49a8-92bc-53003d7237a3"
		permissions := []string{"admin"}

		req.Header.Set("Authorization", generateToken(id, permissions))

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

		id := "10950f72-29ec-49a8-92bc-53003d7237a3"
		permissions := []string{"admin"}

		req.AddCookie(createAccessTokenCookie(id, permissions))

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(ctr.GetProfile)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusForbidden, rr.Code)
	})

	t.Run("should return internal server error when usecase throws error", func(t *testing.T) {
		req, err := createJsonRequest()
		assert.NoError(t, err)

		ctr := handler.NewProfileHandler(&mockGetProfileUseCase{shouldReturnError: true}, nil, nil)

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
		ctx := context.WithValue(req.Context(), cm.RequestIDKey, "test-trace-id")
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

		ph := handler.NewProfileHandler(&mockGetProfileUseCase{}, &mockUpdateUserUseCase{}, nil)
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

		ph := handler.NewProfileHandler(&mockGetProfileUseCase{}, &mockUpdateUserUseCase{}, nil)
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

		ph := handler.NewProfileHandler(&mockGetProfileUseCase{}, &mockUpdateUserUseCase{}, nil)
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

		ph := handler.NewProfileHandler(&mockGetProfileUseCase{}, &mockUpdateUserUseCase{}, nil)
		handler := http.HandlerFunc(ph.UpdateProfile)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("should return FORBIDDEN when user is inactive", func(t *testing.T) {
		req := createJsonRequest(`{
			"username": "coolusername"
		  }`)

		id := "10950f72-29ec-49a8-92bc-53003d7237a3"
		permissions := []string{"admin"}

		req.AddCookie(createAccessTokenCookie(id, permissions))

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

		id := "10950f72-29ec-49a8-92bc-53003d7237a3"
		permissions := []string{"admin"}

		req.AddCookie(createAccessTokenCookie(id, permissions))

		rr := httptest.NewRecorder()

		ph := handler.NewProfileHandler(&mockGetProfileUseCase{}, &mockUpdateUserUseCase{shouldReturnConflictError: true}, nil)
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

		ph := handler.NewProfileHandler(&mockGetProfileUseCase{}, &mockUpdateUserUseCase{shouldReturnError: true}, nil)
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
		req = req.WithContext(ctx)
		return req
	}

	t.Run("should return OK when token is valid", func(t *testing.T) {
		req := createJsonRequest()

		id := "10950f72-29ec-49a8-92bc-53003d7237a3"
		permissions := []string{"admin"}

		req.Header.Set("Authorization", generateToken(id, permissions))

		rr := httptest.NewRecorder()

		ph := handler.NewProfileHandler(&mockGetProfileUseCase{}, nil, &mockDeleteUserUseCase{})
		handler := http.HandlerFunc(ph.DeleteProfile)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNoContent, rr.Code)
	})

	t.Run("should return UNAUTHORIZED when token header is not valid", func(t *testing.T) {
		req := createJsonRequest()

		req.Header.Set("Authorization", "Bearer BadToken")

		rr := httptest.NewRecorder()

		ph := handler.NewProfileHandler(&mockGetProfileUseCase{}, nil, &mockDeleteUserUseCase{})
		handler := http.HandlerFunc(ph.DeleteProfile)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("should return UNAUTHORIZED when token cookie is not valid", func(t *testing.T) {
		req := createJsonRequest()

		req.AddCookie(createBadTokenCookie())

		rr := httptest.NewRecorder()

		ph := handler.NewProfileHandler(&mockGetProfileUseCase{}, nil, &mockDeleteUserUseCase{})
		handler := http.HandlerFunc(ph.DeleteProfile)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("should return INTERNAL SERVER ERROR when usecase returns an error", func(t *testing.T) {
		req := createJsonRequest()

		id := "10950f72-29ec-49a8-92bc-53003d7237a3"
		permissions := []string{"admin"}

		req.AddCookie(createAccessTokenCookie(id, permissions))

		rr := httptest.NewRecorder()

		ph := handler.NewProfileHandler(&mockGetProfileUseCase{}, nil, &mockDeleteUserUseCase{shouldReturnError: true})
		handler := http.HandlerFunc(ph.DeleteProfile)
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
	accountInactive   bool
}

func (u *mockGetProfileUseCase) Execute(ctx context.Context, userId string) (user *model.User, error_ error) {
	if u.shouldReturnError {
		return nil, errors.New("generic error")
	}

	usr := &model.User{
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
	}

	if u.accountInactive {
		usr.IsEnabled = false
	}

	return usr, nil
}

type mockUpdateUserUseCase struct {
	shouldReturnError                bool
	shouldReturnConflictError        bool
	shouldReturnInactiveAccountError bool
}

func (u *mockUpdateUserUseCase) Execute(ctx context.Context, user *model.User) error {
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
