package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	v1 "raspstore.github.io/users-service/api/v1"
	"raspstore.github.io/users-service/internal/api/handler"
	"raspstore.github.io/users-service/internal/model"
)

func TestGetUserConfigurationWithSuccess(t *testing.T) {
	ucr := &userConfigurationRepositoryMock{}
	uch := handler.NewUserConfigHandler(ucr)

	req, _ := http.NewRequest("GET", "/config", nil)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(uch.GetUserConfigs)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, 200)

	var body v1.UserConfigurationResponse
	json.Unmarshal(rr.Body.Bytes(), &body)
	assert.NotNil(t, body)
	assert.NotNil(t, body.AllowLoginWithEmail)
	assert.NotNil(t, body.AllowPublicUserCreation)
	assert.NotNil(t, body.EnforceMfa)
	assert.NotEmpty(t, body.MinPasswordLength)
	assert.NotNil(t, body.StorageLimit)
}

func TestGetUserConfigurationWithInternalServerError(t *testing.T) {
	ucr := &userConfigurationRepositoryMock{shouldReturnError: true}
	uch := handler.NewUserConfigHandler(ucr)

	req, _ := http.NewRequest("GET", "/config", nil)
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), middleware.RequestIDKey, "test-trace-id")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(uch.GetUserConfigs)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, 500)

	var body v1.ErrorResponse
	json.Unmarshal(rr.Body.Bytes(), &body)
	assert.NotNil(t, body)
	assert.NotEmpty(t, body.Code)
	assert.NotEmpty(t, body.Message)
	assert.NotEmpty(t, body.TraceId)
}

func TestPatchUserConfigurationWithPartialBodyWithSuccess(t *testing.T) {
	ucr := &userConfigurationRepositoryMock{}
	uch := handler.NewUserConfigHandler(ucr)

	requestBody := []byte(`{
		"allowLoginWithEmail": false
	  }`)

	req, _ := http.NewRequest("PATCH", "/config", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(uch.UpdateUserConfigs)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, 200)

	var body v1.UserConfigurationResponse
	json.Unmarshal(rr.Body.Bytes(), &body)
	assert.NotNil(t, body)
	assert.NotNil(t, body)
	assert.NotNil(t, body.AllowLoginWithEmail)
	assert.NotNil(t, body.AllowPublicUserCreation)
	assert.NotNil(t, body.EnforceMfa)
	assert.NotEmpty(t, body.MinPasswordLength)
	assert.NotNil(t, body.StorageLimit)
}

func TestPatchUserConfigurationWithInternalServerError(t *testing.T) {
	ucr := &userConfigurationRepositoryMock{shouldReturnError: true}
	uch := handler.NewUserConfigHandler(ucr)

	requestBody := []byte(`{
		"allowLoginWithEmail": false
	  }`)

	req, _ := http.NewRequest("PATCH", "/config", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), middleware.RequestIDKey, "test-trace-id")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(uch.UpdateUserConfigs)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, 500)

	var body v1.ErrorResponse
	json.Unmarshal(rr.Body.Bytes(), &body)
	assert.NotNil(t, body)
	assert.NotEmpty(t, body.Code)
	assert.NotEmpty(t, body.Message)
	assert.NotEmpty(t, body.TraceId)
}

type userConfigurationRepositoryMock struct {
	shouldReturnError          bool
	notAllowPublicUserCreation bool
}

func (ucf *userConfigurationRepositoryMock) Find() (usersConfig *model.UserConfiguration, err error) {
	if ucf.shouldReturnError {
		return nil, mongo.ErrClientDisconnected
	}

	return &model.UserConfiguration{
		StorageLimit:            "3G",
		MinPasswordLength:       8,
		AllowPublicUserCreation: !ucf.notAllowPublicUserCreation,
		AllowLoginWithEmail:     false,
		EnforceMfa:              false,
	}, nil
}

func (ucf *userConfigurationRepositoryMock) Update(usersConfig *model.UserConfiguration) error {
	if ucf.shouldReturnError {
		return mongo.ErrClientDisconnected
	}

	return nil
}
