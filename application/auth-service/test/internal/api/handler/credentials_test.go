package handler_test

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	v1 "raspstore.github.io/auth-service/api/v1"
	"raspstore.github.io/auth-service/internal"
	"raspstore.github.io/auth-service/internal/api/handler"
)

func init() {
	err := godotenv.Load("../../../.env.test")

	if err != nil {
		log.Panicln(err.Error())
	}
}

func TestLoginSuccess(t *testing.T) {
	ctr := handler.NewCredentialsHandler(&mockLoginService{})

	reqBody := []byte(`{"mfaToken": ""}`)
	req, err := http.NewRequest("POST", "/auth-service/login", bytes.NewBuffer(reqBody))

	assert.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Basic "+basic("testuser_ok", "testpassword"))

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ctr.Login)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var res v1.LoginResponse

	err = json.NewDecoder(rr.Body).Decode(&res)

	assert.NoError(t, err)

	assert.NotEmpty(t, res.AccessToken)
	assert.NotEmpty(t, res.RefreshToken)
}

func TestLoginFail(t *testing.T) {
	ctr := handler.NewCredentialsHandler(&mockLoginService{})

	reqBody := []byte(`{"mfaToken": ""}`)
	req, err := http.NewRequest("POST", "/auth-service/login", bytes.NewBuffer(reqBody))

	assert.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Basic "+basic("testuser", "testpassword"))

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ctr.Login)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func basic(username string, password string) string {
	header := username + ":" + password

	return base64.StdEncoding.EncodeToString([]byte(header))
}

type mockLoginService struct{}

func (m *mockLoginService) AuthenticateCredentials(username string, rawPassword string, mfaToken string) (accessToken string, refreshToken string, err error) {
	if strings.HasSuffix(username, "_ok") {
		return "mock_access_token", "mock_refresh_token", nil
	}

	return "", "", internal.ErrIncorrectCredentials
}
