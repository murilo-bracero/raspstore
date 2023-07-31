package handler_test

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	cm "github.com/go-chi/chi/v5/middleware"
	"github.com/murilo-bracero/raspstore/idp/internal"
	"github.com/murilo-bracero/raspstore/idp/internal/api/handler"
	"github.com/murilo-bracero/raspstore/idp/internal/api/middleware"
	"github.com/murilo-bracero/raspstore/idp/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	createJsonRequest := func(body string) (*http.Request, error) {
		reqBody := []byte(body)
		req, err := http.NewRequest("POST", "/idp/v1/admin/users", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		ctx := context.WithValue(req.Context(), cm.RequestIDKey, "test-trace-id")
		ctx = context.WithValue(ctx, middleware.UserClaimsCtxKey, &model.UserClaims{})
		req = req.WithContext(ctx)
		return req, nil
	}

	t.Run("happy path - should return created when payload is valid", func(t *testing.T) {
		req, err := createJsonRequest(`{
			"username": "user6",
			"password": "password123",
			"roles": ["admin"]
		}`)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()

		ah := handler.NewAdminHandler(&mockCreateUserUseCase{}, nil, nil, nil, nil)
		handler := http.HandlerFunc(ah.CreateUser)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("bad request - should return BAD REQUEST when payload is invalid", func(t *testing.T) {
		req, err := createJsonRequest(`{
			"username": "user6",
			"password": "password123"
		}`)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()

		ah := handler.NewAdminHandler(&mockCreateUserUseCase{}, nil, nil, nil, nil)
		handler := http.HandlerFunc(ah.CreateUser)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("bad request - should return BAD REQUEST when conflict", func(t *testing.T) {
		req, err := createJsonRequest(`{
			"username": "user6",
			"password": "password123",
			"roles": ["admin"]
		}`)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()

		ah := handler.NewAdminHandler(&mockCreateUserUseCase{shouldReturnConflict: true}, nil, nil, nil, nil)
		handler := http.HandlerFunc(ah.CreateUser)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("internal server error - should return INTERNAL SERVER ERROR when usecase returns an error", func(t *testing.T) {
		req, err := createJsonRequest(`{
			"username": "user6",
			"password": "password123",
			"roles": ["admin"]
		}`)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()

		ah := handler.NewAdminHandler(&mockCreateUserUseCase{shouldReturnErr: true}, nil, nil, nil, nil)
		handler := http.HandlerFunc(ah.CreateUser)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}

func TestListUser(t *testing.T) {
	createJsonRequest := func() (*http.Request, error) {
		req, err := http.NewRequest("GET", "/idp/v1/admin/users", nil)
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		ctx := context.WithValue(req.Context(), cm.RequestIDKey, "test-trace-id")
		ctx = context.WithValue(ctx, middleware.UserClaimsCtxKey, &model.UserClaims{})
		req = req.WithContext(ctx)
		return req, nil
	}

	t.Run("happy path - should return OK", func(t *testing.T) {
		req, err := createJsonRequest()
		assert.NoError(t, err)

		rr := httptest.NewRecorder()

		ah := handler.NewAdminHandler(nil, nil, nil, &mockListUsersUseCase{}, nil)
		handler := http.HandlerFunc(ah.ListUsers)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("internal server error - should return INTERNAL SERVER ERROR when usecase returns an error", func(t *testing.T) {
		req, err := createJsonRequest()
		assert.NoError(t, err)

		rr := httptest.NewRecorder()

		ah := handler.NewAdminHandler(nil, nil, nil, &mockListUsersUseCase{shouldReturnErr: true}, nil)
		handler := http.HandlerFunc(ah.ListUsers)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}

func TestGetUserById(t *testing.T) {
	createJsonRequest := func() (*http.Request, error) {
		req, err := http.NewRequest("GET", "/idp/v1/admin/users/a9087a4c-2477-4928-8f22-949770a58f8a", nil)
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		ctx := context.WithValue(req.Context(), cm.RequestIDKey, "test-trace-id")
		ctx = context.WithValue(ctx, middleware.UserClaimsCtxKey, &model.UserClaims{})
		req = req.WithContext(ctx)
		return req, nil
	}

	t.Run("happy path - should return OK", func(t *testing.T) {
		req, err := createJsonRequest()
		assert.NoError(t, err)

		rr := httptest.NewRecorder()

		ah := handler.NewAdminHandler(nil, &mockGetUserUseCase{}, nil, nil, nil)
		handler := http.HandlerFunc(ah.GetUserById)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("not found - should return NOT FOUND", func(t *testing.T) {
		req, err := createJsonRequest()
		assert.NoError(t, err)

		rr := httptest.NewRecorder()

		ah := handler.NewAdminHandler(nil, &mockGetUserUseCase{shouldReturnNotFound: true}, nil, nil, nil)
		handler := http.HandlerFunc(ah.GetUserById)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNotFound, rr.Code)
	})

	t.Run("internal server error - should return INTERNAL SERVER ERROR", func(t *testing.T) {
		req, err := createJsonRequest()
		assert.NoError(t, err)

		rr := httptest.NewRecorder()

		ah := handler.NewAdminHandler(nil, &mockGetUserUseCase{shouldReturnErr: true}, nil, nil, nil)
		handler := http.HandlerFunc(ah.GetUserById)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}

func TestUpdateUserById(t *testing.T) {
	createJsonRequest := func(body string) (*http.Request, error) {
		reqBody := []byte(body)
		req, err := http.NewRequest("PUT", "/idp/v1/admin/users/a9087a4c-2477-4928-8f22-949770a58f8a", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		ctx := context.WithValue(req.Context(), cm.RequestIDKey, "test-trace-id")
		ctx = context.WithValue(ctx, middleware.UserClaimsCtxKey, &model.UserClaims{})
		req = req.WithContext(ctx)
		return req, nil
	}

	t.Run("happy path - should return OK", func(t *testing.T) {
		req, err := createJsonRequest(`{
			"enabled": false,
			"username": "coolusername2",
			"mfaEnabled": false,
			"roles": ["user"]
		}`)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()

		ah := handler.NewAdminHandler(nil, nil, nil, nil, &mockUpdateUserUseCase{})
		handler := http.HandlerFunc(ah.UpdateUserById)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("bad request - should return BAD REQUEST when payload is invalid", func(t *testing.T) {
		req, err := createJsonRequest(`{
			"enabled": false,
			"username": "coolusername2",
			"mfaEnabled": false
		}`)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()

		ah := handler.NewAdminHandler(nil, nil, nil, nil, &mockUpdateUserUseCase{})
		handler := http.HandlerFunc(ah.UpdateUserById)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("not found - should return NOT FOUND", func(t *testing.T) {
		req, err := createJsonRequest(`{
			"enabled": false,
			"username": "coolusername2",
			"mfaEnabled": false,
			"roles": ["user"]
		}`)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()

		ah := handler.NewAdminHandler(nil, nil, nil, nil, &mockUpdateUserUseCase{shouldReturnNotFound: true})
		handler := http.HandlerFunc(ah.UpdateUserById)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNotFound, rr.Code)
	})

	t.Run("conflict - should return CONFLICT", func(t *testing.T) {
		req, err := createJsonRequest(`{
			"enabled": false,
			"username": "coolusername2",
			"mfaEnabled": false,
			"roles": ["user"]
		}`)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()

		ah := handler.NewAdminHandler(nil, nil, nil, nil, &mockUpdateUserUseCase{shouldReturnConflict: true})
		handler := http.HandlerFunc(ah.UpdateUserById)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusConflict, rr.Code)
	})

	t.Run("internal server error - should return INTERNAL SERVER ERROR", func(t *testing.T) {
		req, err := createJsonRequest(`{
			"enabled": false,
			"username": "coolusername2",
			"mfaEnabled": false,
			"roles": ["user"]
		}`)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()

		ah := handler.NewAdminHandler(nil, nil, nil, nil, &mockUpdateUserUseCase{shouldReturnErr: true})
		handler := http.HandlerFunc(ah.UpdateUserById)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}

func TestDeleteUser(t *testing.T) {
	createJsonRequest := func() (*http.Request, error) {
		req, err := http.NewRequest("DELETE", "/idp/v1/admin/users/a9087a4c-2477-4928-8f22-949770a58f8a", nil)
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		ctx := context.WithValue(req.Context(), cm.RequestIDKey, "test-trace-id")
		ctx = context.WithValue(ctx, middleware.UserClaimsCtxKey, &model.UserClaims{})
		req = req.WithContext(ctx)
		return req, nil
	}

	t.Run("happy path - should return NO CONTENT", func(t *testing.T) {
		req, err := createJsonRequest()
		assert.NoError(t, err)

		rr := httptest.NewRecorder()

		ah := handler.NewAdminHandler(nil, nil, &mockDeleteUserUseCase{}, nil, nil)
		handler := http.HandlerFunc(ah.DeleteUser)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNoContent, rr.Code)
	})

	t.Run("internal server error - should return INTERNAL SERVER ERROR", func(t *testing.T) {
		req, err := createJsonRequest()
		assert.NoError(t, err)

		rr := httptest.NewRecorder()

		ah := handler.NewAdminHandler(nil, nil, &mockDeleteUserUseCase{shouldReturnError: true}, nil, nil)
		handler := http.HandlerFunc(ah.DeleteUser)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}

type mockCreateUserUseCase struct {
	shouldReturnErr      bool
	shouldReturnConflict bool
}

func (u *mockCreateUserUseCase) Execute(ctx context.Context, user *model.User) error {
	if u.shouldReturnErr {
		return errors.New("generic error")
	}

	if u.shouldReturnConflict {
		return internal.ErrConflict
	}

	return nil
}

type mockListUsersUseCase struct {
	shouldReturnErr bool
}

func (u *mockListUsersUseCase) Execute(ctx context.Context, page int, size int, username string, enabled *bool) (userPage *model.UserPage, error_ error) {
	if u.shouldReturnErr {
		return nil, errors.New("generic error")
	}

	return &model.UserPage{}, nil
}

type mockGetUserUseCase struct {
	shouldReturnErr      bool
	shouldReturnNotFound bool
}

func (u *mockGetUserUseCase) Execute(ctx context.Context, userId string) (user *model.User, error_ error) {
	if u.shouldReturnErr {
		return nil, errors.New("generic error")
	}

	if u.shouldReturnNotFound {
		return nil, internal.ErrUserNotFound
	}

	return &model.User{}, nil
}

type mockUpdateUserUseCase struct {
	shouldReturnErr      bool
	shouldReturnNotFound bool
	shouldReturnConflict bool
}

func (u *mockUpdateUserUseCase) Execute(ctx context.Context, user *model.User) error {
	if u.shouldReturnErr {
		return errors.New("generic error")
	}

	if u.shouldReturnNotFound {
		return internal.ErrUserNotFound
	}

	if u.shouldReturnConflict {
		return internal.ErrConflict
	}

	return nil
}
