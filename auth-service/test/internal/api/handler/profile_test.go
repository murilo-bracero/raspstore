package handler_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	v1 "github.com/murilo-bracero/raspstore/auth-service/api/v1"
	"github.com/murilo-bracero/raspstore/auth-service/internal/api/handler"
	"github.com/murilo-bracero/raspstore/auth-service/internal/model"
	"github.com/murilo-bracero/raspstore/auth-service/internal/token"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestGetProfileSuccessWithTokenInHeader(t *testing.T) {

	ctr := handler.NewProfileHandler(&mockGetProfileUseCase{})

	req, err := http.NewRequest("GET", "/idp/v1/profile", nil)

	assert.NoError(t, err)

	id := "10950f72-29ec-49a8-92bc-53003d7237a3"
	permissions := []string{"admin"}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", generateToken(id, permissions))

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ctr.GetProfile)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var res v1.UserRepresentation

	err = json.NewDecoder(rr.Body).Decode(&res)

	assert.NoError(t, err)

	assert.Equal(t, "c223a9f5-7174-4102-aacc-73f03954dde8", res.UserID)
	assert.Equal(t, "cool_username", res.Username)
	assert.Equal(t, "cool.email@email.com", res.Email)
	assert.Equal(t, "123-555-456", res.PhoneNumber)
	assert.Equal(t, true, res.IsMfaEnabled)
	assert.Equal(t, false, res.IsMfaVerified)

	assert.NotEmpty(t, res.CreatedAt)
	assert.NotEmpty(t, res.UpdatedAt)
}

func TestGetProfileSuccessWithTokenInCookie(t *testing.T) {

	ctr := handler.NewProfileHandler(&mockGetProfileUseCase{})

	req, err := http.NewRequest("GET", "/idp/v1/profile", nil)

	assert.NoError(t, err)

	id := "10950f72-29ec-49a8-92bc-53003d7237a3"
	permissions := []string{"admin"}

	req.AddCookie(createAccessTokenCookie(id, permissions))

	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ctr.GetProfile)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var res v1.UserRepresentation

	err = json.NewDecoder(rr.Body).Decode(&res)

	assert.NoError(t, err)

	assert.Equal(t, "c223a9f5-7174-4102-aacc-73f03954dde8", res.UserID)
	assert.Equal(t, "cool_username", res.Username)
	assert.Equal(t, "cool.email@email.com", res.Email)
	assert.Equal(t, "123-555-456", res.PhoneNumber)
	assert.Equal(t, true, res.IsMfaEnabled)
	assert.Equal(t, false, res.IsMfaVerified)

	assert.NotEmpty(t, res.CreatedAt)
	assert.NotEmpty(t, res.UpdatedAt)
}

func TestGetProfileSuccessWithoutToken(t *testing.T) {

	ctr := handler.NewProfileHandler(&mockGetProfileUseCase{})

	req, err := http.NewRequest("GET", "/idp/v1/profile", nil)

	assert.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ctr.GetProfile)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestGetProfileUnauthorizedWithMalformedTokenInCookie(t *testing.T) {

	ctr := handler.NewProfileHandler(&mockGetProfileUseCase{})

	req, err := http.NewRequest("GET", "/idp/v1/profile", nil)

	assert.NoError(t, err)

	req.AddCookie(createBadTokenCookie())

	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ctr.GetProfile)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestGetProfileUnauthorizedWithMalformedTokenInHeader(t *testing.T) {

	ctr := handler.NewProfileHandler(&mockGetProfileUseCase{})

	req, err := http.NewRequest("GET", "/idp/v1/profile", nil)

	assert.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer badToken")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ctr.GetProfile)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestGetProfileInternalServerError(t *testing.T) {

	ctr := handler.NewProfileHandler(&mockGetProfileUseCase{shouldReturnError: true})

	req, err := http.NewRequest("GET", "/idp/v1/profile", nil)

	assert.NoError(t, err)

	id := "10950f72-29ec-49a8-92bc-53003d7237a3"
	permissions := []string{"admin"}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", generateToken(id, permissions))

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ctr.GetProfile)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
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
		Email:         "cool.email@email.com",
		IsEnabled:     true,
		PhoneNumber:   "123-555-456",
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
