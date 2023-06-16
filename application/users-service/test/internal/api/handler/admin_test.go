package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	v1 "raspstore.github.io/users-service/api/v1"
	"raspstore.github.io/users-service/internal/api/handler"
)

func TestAdminSaveUserWithSuccess(t *testing.T) {
	ur := &userServiceMock{}
	ctr := handler.NewAdminUserHandler(ur)

	random := uuid.NewString()
	reqBody := []byte(fmt.Sprintf(`{
		"username": "%s",
		"email": "%s@test.com",
		"password": "%s_super-secret-password"
	  }`, random, random, random))
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), middleware.RequestIDKey, "test-trace-id")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ctr.CreateUser)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, 201)

	var usrRes v1.AdminUserResponse
	json.Unmarshal(rr.Body.Bytes(), &usrRes)
	assert.NotNil(t, usrRes)
	assert.NotEmpty(t, usrRes.UserId)
	assert.NotEmpty(t, usrRes.Email)
	assert.NotEmpty(t, usrRes.Username)
	assert.NotEmpty(t, usrRes.CreatedAt)
	assert.NotEmpty(t, usrRes.UpdatedAt)
}

func TestAdminSaveUserWithPermissionsWithSuccess(t *testing.T) {
	ur := &userServiceMock{}
	ctr := handler.NewAdminUserHandler(ur)

	random := uuid.NewString()
	reqBody := []byte(fmt.Sprintf(`{
		"username": "%s",
		"email": "%s@test.com",
		"password": "%s_super-secret-password",
		"permissions": ["admin"]
	  }`, random, random, random))
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), middleware.RequestIDKey, "test-trace-id")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ctr.CreateUser)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, 201)

	var usrRes v1.AdminUserResponse
	json.Unmarshal(rr.Body.Bytes(), &usrRes)
	assert.NotNil(t, usrRes)
	assert.NotEmpty(t, usrRes.UserId)
	assert.NotEmpty(t, usrRes.Email)
	assert.NotEmpty(t, usrRes.Username)
	assert.NotNil(t, usrRes.Permissions)
	assert.NotEmpty(t, usrRes.CreatedAt)
	assert.NotEmpty(t, usrRes.UpdatedAt)
}
