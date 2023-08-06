package handler_test

import (
	"bytes"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	v1 "github.com/murilo-bracero/raspstore/idp/api/v1"
	"github.com/murilo-bracero/raspstore/idp/internal"
	"github.com/murilo-bracero/raspstore/idp/internal/api/handler"
	"github.com/murilo-bracero/raspstore/idp/internal/api/utils"
	"github.com/murilo-bracero/raspstore/idp/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	createReq := func(responseType string, username string, password string) *http.Request {
		data := url.Values{}
		data.Set("response_type", responseType)
		req, err := http.NewRequest("POST", "/idp/login", bytes.NewBufferString(data.Encode()))

		assert.NoError(t, err)

		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("Authorization", "Basic "+basic(username, password))
		return req
	}

	t.Run("Should return success with json tokens when username/password is correct and response_type is code", func(t *testing.T) {
		ctr := handler.NewLoginHandler(&mockLoginUseCase{})

		req := createReq("code", "testuser_ok", "password123")
		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(ctr.Login)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var res v1.LoginResponse

		err := utils.ParseBody(rr.Body, &res)
		assert.NoError(t, err)

		assert.NotEmpty(t, res.AccessToken)
		assert.NotEmpty(t, res.RefreshToken)
	})

	t.Run("Should return success with cookie tokens when username/password is correct and response_type is token", func(t *testing.T) {
		ctr := handler.NewLoginHandler(&mockLoginUseCase{})

		req := createReq("token", "testuser_ok", "password123")
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
	})

	t.Run("Should return UNAUTHORIZED when username/password is wrong", func(t *testing.T) {
		ctr := handler.NewLoginHandler(&mockLoginUseCase{})

		req := createReq("token", "testuser", "password123")
		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(ctr.Login)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})
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
