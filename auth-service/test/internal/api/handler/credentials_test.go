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
	"time"

	"github.com/joho/godotenv"
	v1 "github.com/murilo-bracero/raspstore/auth-service/api/v1"
	"github.com/murilo-bracero/raspstore/auth-service/internal"
	"github.com/murilo-bracero/raspstore/auth-service/internal/api/handler"
	"github.com/murilo-bracero/raspstore/auth-service/internal/model"
	"github.com/stretchr/testify/assert"
)

func init() {
	err := godotenv.Load("../../../.env.test")

	if err != nil {
		log.Panicln(err.Error())
	}
}

func TestLoginSuccess(t *testing.T) {
	ctr := handler.NewLoginHandler(&mockLoginUseCase{})

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
	ctr := handler.NewLoginHandler(&mockLoginUseCase{})

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

type mockLoginUseCase struct{}

func (m *mockLoginUseCase) AuthenticateCredentials(username string, rawPassword string, mfaToken string) (tokenCredentials *model.TokenCredentials, err error) {
	if strings.HasSuffix(username, "_ok") {
		return &model.TokenCredentials{
			AccessToken:  "mock_access_token",
			RefreshToken: "mock_refresh_token",
			ExpirestAt:   time.Now().Add(24 * time.Hour),
		}, nil
	}

	return nil, internal.ErrIncorrectCredentials
}
