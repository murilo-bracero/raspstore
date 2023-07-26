package handler_test

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
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

func TestLoginSuccessWithCodeResponseType(t *testing.T) {
	ctr := handler.NewLoginHandler(&mockLoginUseCase{})

	data := url.Values{}
	data.Set("response_type", "code")
	req, err := http.NewRequest("POST", "/auth-service/login", bytes.NewBufferString(data.Encode()))

	assert.NoError(t, err)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
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

func TestLoginSuccessWithTokenResponseType(t *testing.T) {
	ctr := handler.NewLoginHandler(&mockLoginUseCase{})

	data := url.Values{}
	data.Set("response_type", "token")
	req, err := http.NewRequest("POST", "/auth-service/login", bytes.NewBufferString(data.Encode()))

	assert.NoError(t, err)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Basic "+basic("testuser_ok", "testpassword"))

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ctr.Login)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	accessTokenCookie := findCookieByName(rr.Result().Cookies(), "access_token")
	assert.NotNil(t, accessTokenCookie)

	refreshTokenCookie := findCookieByName(rr.Result().Cookies(), "refresh_token")
	assert.NotNil(t, refreshTokenCookie)

	assert.NotEmpty(t, accessTokenCookie)
	assert.NotEmpty(t, refreshTokenCookie)
}

func TestLoginFail(t *testing.T) {
	ctr := handler.NewLoginHandler(&mockLoginUseCase{})

	data := url.Values{}
	data.Set("response_type", "code")
	req, err := http.NewRequest("POST", "/auth-service/login", bytes.NewBufferString(data.Encode()))

	assert.NoError(t, err)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
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

func findCookieByName(cookies []*http.Cookie, name string) *http.Cookie {
	for _, cookie := range cookies {
		if cookie.Name == name {
			return cookie
		}
	}

	return nil
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
