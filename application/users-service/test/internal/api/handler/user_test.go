package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"github.com/murilo-bracero/raspstore/commons/pkg/middleware"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	v1 "raspstore.github.io/users-service/api/v1"
	"raspstore.github.io/users-service/internal"
	"raspstore.github.io/users-service/internal/api/handler"
	"raspstore.github.io/users-service/internal/model"
)

func TestGetUserWithSuccess(t *testing.T) {
	ur := &userServiceMock{}
	ucr := &userConfigurationRepositoryMock{}
	ctr := handler.NewUserHandler(ur, ucr)

	random := uuid.NewString()
	req, _ := http.NewRequest("GET", fmt.Sprintf("/users/%s", random), nil)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ctr.GetUser)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, 200)

	var usrRes v1.UserResponse
	json.Unmarshal(rr.Body.Bytes(), &usrRes)
	assert.NotNil(t, usrRes)
	assert.NotEmpty(t, usrRes.UserId)
	assert.NotEmpty(t, usrRes.Email)
	assert.NotEmpty(t, usrRes.Username)
	assert.NotEmpty(t, usrRes.CreatedAt)
	assert.NotEmpty(t, usrRes.UpdatedAt)
}

func TestGetUserWithNotFoundError(t *testing.T) {
	ur := &userServiceMock{throwNotFound: true}
	ucr := &userConfigurationRepositoryMock{}
	ctr := handler.NewUserHandler(ur, ucr)

	random := uuid.NewString()
	req, _ := http.NewRequest("GET", fmt.Sprintf("/users/%s", random), nil)
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), chiMiddleware.RequestIDKey, "test-trace-id")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ctr.GetUser)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestGetUserWithInternalServerError(t *testing.T) {
	ur := &userServiceMock{throwInternalError: true}
	ucr := &userConfigurationRepositoryMock{}
	ctr := handler.NewUserHandler(ur, ucr)

	random := uuid.NewString()
	req, _ := http.NewRequest("GET", fmt.Sprintf("/users/%s", random), nil)
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), chiMiddleware.RequestIDKey, "test-trace-id")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ctr.GetUser)
	handler.ServeHTTP(rr, req)

	var errResponse v1.ErrorResponse
	json.Unmarshal(rr.Body.Bytes(), &errResponse)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.NotEmpty(t, errResponse.Code)
	assert.NotEmpty(t, errResponse.Message)
	assert.NotEmpty(t, errResponse.TraceId)
}

func TestSaveUserWithSuccess(t *testing.T) {
	ur := &userServiceMock{}
	ucr := &userConfigurationRepositoryMock{}
	ctr := handler.NewUserHandler(ur, ucr)

	random := uuid.NewString()
	reqBody := []byte(fmt.Sprintf(`{
		"username": "%s",
		"email": "%s@test.com",
		"password": "%s_super-secret-password"
	  }`, random, random, random))
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), chiMiddleware.RequestIDKey, "test-trace-id")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ctr.CreateUser)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, 201)

	var usrRes v1.UserResponse
	json.Unmarshal(rr.Body.Bytes(), &usrRes)
	assert.NotNil(t, usrRes)
	assert.NotEmpty(t, usrRes.UserId)
	assert.NotEmpty(t, usrRes.Email)
	assert.NotEmpty(t, usrRes.Username)
	assert.NotEmpty(t, usrRes.CreatedAt)
	assert.NotEmpty(t, usrRes.UpdatedAt)
}

func TestSaveUserWithErrorWhenPublicUserCreationNotAllowed(t *testing.T) {
	ur := &userServiceMock{}
	ucr := &userConfigurationRepositoryMock{notAllowPublicUserCreation: true}
	ctr := handler.NewUserHandler(ur, ucr)

	random := uuid.NewString()
	reqBody := []byte(fmt.Sprintf(`{
		"username": "%s",
		"email": "%s@test.com",
		"password": "%s_super-secret-password"
	  }`, random, random, random))
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), chiMiddleware.RequestIDKey, "test-trace-id")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ctr.CreateUser)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, 400)

	var usrRes v1.ErrorResponse
	json.Unmarshal(rr.Body.Bytes(), &usrRes)
	assert.NotNil(t, usrRes)
	assert.NotEmpty(t, usrRes.Code)
	assert.NotEmpty(t, usrRes.Message)
	assert.NotEmpty(t, usrRes.TraceId)
}

func TestSaveUserWithInvalidPayload(t *testing.T) {
	ur := &userServiceMock{}
	ucr := &userConfigurationRepositoryMock{}
	ctr := handler.NewUserHandler(ur, ucr)

	random := uuid.NewString()
	reqBody := []byte(fmt.Sprintf(`{
		"username": "%s",
		"password": "%s_super-secret-password"
	  }`, random, random))
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), chiMiddleware.RequestIDKey, "test-trace-id")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ctr.CreateUser)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, 400)

	var errRes v1.ErrorResponse
	json.Unmarshal(rr.Body.Bytes(), &errRes)
	assert.NotEmpty(t, errRes.Code)
	assert.NotEmpty(t, errRes.Message)
	assert.NotEmpty(t, errRes.TraceId)
}

func TestSaveUserWithAlreadyExistedEmail(t *testing.T) {
	ur := &userServiceMock{throwAlreadyExists: true}
	ucr := &userConfigurationRepositoryMock{}
	ctr := handler.NewUserHandler(ur, ucr)

	random := uuid.NewString()
	reqBody := []byte(fmt.Sprintf(`{
		"username": "%s",
		"email": "existing@test.com",
		"password": "%s_super-secret-password"
	  }`, random, random))
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), chiMiddleware.RequestIDKey, "test-trace-id")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ctr.CreateUser)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, 400)

	var errRes v1.ErrorResponse
	json.Unmarshal(rr.Body.Bytes(), &errRes)
	assert.NotEmpty(t, errRes.Code)
	assert.NotEmpty(t, errRes.Message)
	assert.NotEmpty(t, errRes.TraceId)
}

func TestSaveUserWithInternalServerError(t *testing.T) {
	ur := &userServiceMock{throwInternalError: true}
	ucr := &userConfigurationRepositoryMock{}
	ctr := handler.NewUserHandler(ur, ucr)

	random := uuid.NewString()
	reqBody := []byte(fmt.Sprintf(`{
		"username": "%s",
		"email": "%s@test.com",
		"password": "%s_super-secret-password"
	  }`, random, random, random))
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), chiMiddleware.RequestIDKey, "test-trace-id")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ctr.CreateUser)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	var errRes v1.ErrorResponse
	json.Unmarshal(rr.Body.Bytes(), &errRes)
	assert.NotEmpty(t, errRes.Code)
	assert.NotEmpty(t, errRes.Message)
	assert.NotEmpty(t, errRes.TraceId)
}

func makeRequest(method string, route string, body []byte) *http.Request {
	var req *http.Request

	if body != nil {
		req, _ = http.NewRequest(method, route, bytes.NewBuffer(body))
	} else {
		req, _ = http.NewRequest(method, route, nil)
	}

	req.Header.Set("Content-Type", "application/json")
	rctx := chi.NewRouteContext()

	if strings.Contains(route, "{id}") {
		rctx.URLParams.Add("id", uuid.NewString())
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	}

	req = req.WithContext(context.WithValue(req.Context(), chiMiddleware.RequestIDKey, "test-trace-id"))
	return req.WithContext(context.WithValue(req.Context(), middleware.UserIdKey, "requester-id"))
}

func TestUpdateUserSuccess(t *testing.T) {
	ur := &userServiceMock{}
	ucr := &userConfigurationRepositoryMock{}
	ctr := handler.NewUserHandler(ur, ucr)

	random := uuid.NewString()
	reqBody := []byte(fmt.Sprintf(`{
		"username": "%s",
		"email": "updated_%s@test.com"
	  }`, random, random))
	req := makeRequest("PUT", "/users/{id}", reqBody)

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ctr.UpdateUser)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var usrRes v1.UserResponse
	json.Unmarshal(rr.Body.Bytes(), &usrRes)
	assert.NotNil(t, usrRes)
	assert.NotEmpty(t, usrRes.UserId)
	assert.NotEmpty(t, usrRes.Email)
	assert.NotEmpty(t, usrRes.Username)
	assert.NotEmpty(t, usrRes.CreatedAt)
	assert.NotEmpty(t, usrRes.UpdatedAt)
}

func TestUpdateUserWithAlreadyExistedNewEmail(t *testing.T) {
	ur := &userServiceMock{throwAlreadyExists: true}
	ucr := &userConfigurationRepositoryMock{}
	ctr := handler.NewUserHandler(ur, ucr)

	random := uuid.NewString()
	reqBody := []byte(`{ "email": "existed@test.com"}`)
	req, _ := http.NewRequest("PUT", fmt.Sprintf("/users/%s", random), bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), chiMiddleware.RequestIDKey, "test-trace-id")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ctr.UpdateUser)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusConflict, rr.Code)
}

func TestUpdateUserInternalServerError(t *testing.T) {
	ur := &userServiceMock{throwInternalError: true}
	ucr := &userConfigurationRepositoryMock{}
	ctr := handler.NewUserHandler(ur, ucr)

	random := uuid.NewString()
	reqBody := []byte(fmt.Sprintf(`{
		"username": "%s",
		"email": "new_%s@test.com",
		"isEnabled": true,
		"password": "%s_super-secret-password"
	  }`, random, random, random))
	req, _ := http.NewRequest("PUT", fmt.Sprintf("/users/%s", random), bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), chiMiddleware.RequestIDKey, "test-trace-id")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ctr.UpdateUser)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	var errRes v1.ErrorResponse
	json.Unmarshal(rr.Body.Bytes(), &errRes)
	assert.NotEmpty(t, errRes.Code)
	assert.NotEmpty(t, errRes.Message)
	assert.NotEmpty(t, errRes.TraceId)
}

func TestListUsers(t *testing.T) {
	ur := &userServiceMock{}
	ucr := &userConfigurationRepositoryMock{}
	ctr := handler.NewUserHandler(ur, ucr)

	req, _ := http.NewRequest("GET", "/users", nil)
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ctr.ListUser)
	handler.ServeHTTP(rr, req)

	var userPage model.UserPage
	json.Unmarshal(rr.Body.Bytes(), &userPage)

	assert.Equal(t, http.StatusOK, rr.Code)

	assert.True(t, len(userPage.Content) > 0)
}

func TestListUsersWithPagination(t *testing.T) {
	totalElements := 10
	ur := &userServiceMock{totalElements: totalElements}
	ucr := &userConfigurationRepositoryMock{}
	ctr := handler.NewUserHandler(ur, ucr)

	size := 2
	page := 0

	req, _ := http.NewRequest("GET", fmt.Sprintf("/users?page=%d&size=%d", page, size), nil)
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ctr.ListUser)
	handler.ServeHTTP(rr, req)

	var userResponseList v1.UserResponseList
	json.Unmarshal(rr.Body.Bytes(), &userResponseList)

	assert.Equal(t, http.StatusOK, rr.Code)

	assert.Equal(t, page, userResponseList.Page)
	assert.Equal(t, size, userResponseList.Size)
	assert.Equal(t, totalElements, userResponseList.TotalElements)
	assert.True(t, strings.HasSuffix(userResponseList.Next, fmt.Sprintf("page=%d&size=%d", page+1, size)))
	assert.Equal(t, size, len(userResponseList.Content))
}

func TestDeleteSuccess(t *testing.T) {
	ur := &userServiceMock{}
	ucr := &userConfigurationRepositoryMock{}
	ctr := handler.NewUserHandler(ur, ucr)

	req := makeRequest("DELETE", "/users/{id}", nil)
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ctr.DeleteUser)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)
}

func TestDeleteInternalServerError(t *testing.T) {
	ur := &userServiceMock{throwInternalError: true}
	ucr := &userConfigurationRepositoryMock{}
	ctr := handler.NewUserHandler(ur, ucr)

	random := uuid.NewString()
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/users/%s", random), nil)
	rr := httptest.NewRecorder()
	ctx := context.WithValue(req.Context(), chiMiddleware.RequestIDKey, "test-trace-id")
	req = req.WithContext(ctx)

	handler := http.HandlerFunc(ctr.DeleteUser)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	var errRes v1.ErrorResponse
	json.Unmarshal(rr.Body.Bytes(), &errRes)
	assert.NotEmpty(t, errRes.Code)
	assert.NotEmpty(t, errRes.Message)
	assert.NotEmpty(t, errRes.TraceId)
}

type userServiceMock struct {
	throwInternalError bool
	throwAlreadyExists bool
	throwNotFound      bool
	totalElements      int
}

func (u *userServiceMock) CreateUser(user *model.User) error {
	if u.throwAlreadyExists {
		return internal.ErrUserAlreadyExists
	}

	if u.throwInternalError {
		return errors.New("generic error")
	}

	return nil
}

func (u *userServiceMock) GetUserById(id string) (*model.User, error) {
	if u.throwInternalError {
		return nil, errors.New("generic error")
	}

	if u.throwNotFound {
		return nil, internal.ErrUserNotFound
	}

	return createRandomUser(id, ""), nil
}

func (u *userServiceMock) RemoveUserById(id string) error {
	if u.throwInternalError {
		return errors.New("generic error")
	}

	return nil
}

func (u *userServiceMock) UpdateUser(user *model.User) (*model.User, error) {
	if u.throwAlreadyExists {
		return nil, internal.ErrEmailOrUsernameInUse
	}

	if u.throwInternalError {
		return nil, errors.New("generic error")
	}

	user.IsEnabled = true
	user.UpdatedAt = time.Now()
	return user, nil
}

func (u *userServiceMock) GetAllUsersByPage(page int, size int) (*model.UserPage, error) {

	users := make([]*model.User, 0)

	for i := 0; i < size; i++ {
		users = append(users, createRandomUser("", ""))
	}

	userPage := &model.UserPage{
		Content: users,
		Count:   u.totalElements,
	}

	return userPage, nil
}

func createRandomUser(id string, email string) *model.User {
	if email == "" {
		email = uuid.NewString() + "@email.com"
	}

	if id == "" {
		id = primitive.NewObjectID().Hex()
	}

	return &model.User{
		UserId:      id,
		Username:    uuid.NewString(),
		Email:       email,
		IsEnabled:   true,
		PhoneNumber: "",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}
