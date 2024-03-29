package controller_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"raspstore.github.io/users-service/api/controller"
	"raspstore.github.io/users-service/api/dto"
	"raspstore.github.io/users-service/model"
)

func TestGetUserWithSuccess(t *testing.T) {
	ur := &userRepositoryMock{}
	ctr := controller.NewUserController(ur)

	random := uuid.NewString()
	req, _ := http.NewRequest("GET", fmt.Sprintf("/users/%s", random), nil)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ctr.GetUser)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, 200)

	var usrRes dto.UserResponse
	json.Unmarshal(rr.Body.Bytes(), &usrRes)
	assert.NotNil(t, usrRes)
	assert.NotEmpty(t, usrRes.UserId)
	assert.NotEmpty(t, usrRes.Email)
	assert.NotEmpty(t, usrRes.Username)
	assert.True(t, usrRes.IsEnabled)
	assert.NotEmpty(t, usrRes.CreatedAt)
	assert.NotEmpty(t, usrRes.UpdatedAt)
}

func TestGetUserWithNotFoundError(t *testing.T) {
	ur := &userRepositoryMock{shouldReturn404: true}
	ctr := controller.NewUserController(ur)

	random := uuid.NewString()
	req, _ := http.NewRequest("GET", fmt.Sprintf("/users/%s", random), nil)
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), middleware.RequestIDKey, "test-trace-id")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ctr.GetUser)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestGetUserWithInternalServerError(t *testing.T) {
	ur := &userRepositoryMock{shouldReturn500: true}
	ctr := controller.NewUserController(ur)

	random := uuid.NewString()
	req, _ := http.NewRequest("GET", fmt.Sprintf("/users/%s", random), nil)
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), middleware.RequestIDKey, "test-trace-id")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ctr.GetUser)
	handler.ServeHTTP(rr, req)

	var errResponse dto.ErrorResponse
	json.Unmarshal(rr.Body.Bytes(), &errResponse)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.NotEmpty(t, errResponse.Code)
	assert.NotEmpty(t, errResponse.Message)
	assert.NotEmpty(t, errResponse.TraceId)
}

func TestSaveUserWithSuccess(t *testing.T) {
	ur := &userRepositoryMock{}
	ctr := controller.NewUserController(ur)

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

	var usrRes dto.UserResponse
	json.Unmarshal(rr.Body.Bytes(), &usrRes)
	assert.NotNil(t, usrRes)
	assert.NotEmpty(t, usrRes.UserId)
	assert.NotEmpty(t, usrRes.Email)
	assert.NotEmpty(t, usrRes.Username)
	assert.True(t, usrRes.IsEnabled)
	assert.NotEmpty(t, usrRes.CreatedAt)
	assert.NotEmpty(t, usrRes.UpdatedAt)
}

func TestSaveUserWithInvalidPayload(t *testing.T) {
	ur := &userRepositoryMock{}
	ctr := controller.NewUserController(ur)

	random := uuid.NewString()
	reqBody := []byte(fmt.Sprintf(`{
		"username": "%s",
		"password": "%s_super-secret-password"
	  }`, random, random))
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), middleware.RequestIDKey, "test-trace-id")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ctr.CreateUser)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, 400)

	var errRes dto.ErrorResponse
	json.Unmarshal(rr.Body.Bytes(), &errRes)
	assert.NotEmpty(t, errRes.Code)
	assert.NotEmpty(t, errRes.Message)
	assert.NotEmpty(t, errRes.TraceId)
}

func TestSaveUserWithAlreadyExistedEmail(t *testing.T) {
	ur := &userRepositoryMock{shouldReturn409: true}
	ctr := controller.NewUserController(ur)

	random := uuid.NewString()
	reqBody := []byte(fmt.Sprintf(`{
		"username": "%s",
		"email": "existing@test.com",
		"password": "%s_super-secret-password"
	  }`, random, random))
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), middleware.RequestIDKey, "test-trace-id")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ctr.CreateUser)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, 400)

	var errRes dto.ErrorResponse
	json.Unmarshal(rr.Body.Bytes(), &errRes)
	assert.NotEmpty(t, errRes.Code)
	assert.NotEmpty(t, errRes.Message)
	assert.NotEmpty(t, errRes.TraceId)
}

func TestSaveUserWithInternalServerError(t *testing.T) {
	ur := &userRepositoryMock{shouldReturn500: true}
	ctr := controller.NewUserController(ur)

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

	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	var errRes dto.ErrorResponse
	json.Unmarshal(rr.Body.Bytes(), &errRes)
	assert.NotEmpty(t, errRes.Code)
	assert.NotEmpty(t, errRes.Message)
	assert.NotEmpty(t, errRes.TraceId)
}

func TestUpdateUserSuccess(t *testing.T) {
	ur := &userRepositoryMock{}
	ctr := controller.NewUserController(ur)

	random := uuid.NewString()
	reqBody := []byte(fmt.Sprintf(`{
		"username": "%s",
		"email": "updated_%s@test.com",
		"isEnabled": true,
		"password": "%s_super-secret-password"
	  }`, random, random, random))
	req, _ := http.NewRequest("PUT", fmt.Sprintf("/users/%s", random), bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ctr.UpdateUser)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var usrRes dto.UserResponse
	json.Unmarshal(rr.Body.Bytes(), &usrRes)
	assert.NotNil(t, usrRes)
	assert.NotEmpty(t, usrRes.UserId)
	assert.NotEmpty(t, usrRes.Email)
	assert.NotEmpty(t, usrRes.Username)
	assert.True(t, usrRes.IsEnabled)
	assert.NotEmpty(t, usrRes.CreatedAt)
	assert.NotEmpty(t, usrRes.UpdatedAt)
}

func TestUpdateUserWithAlreadyExistedNewEmail(t *testing.T) {
	ur := &userRepositoryMock{shouldReturn409: true}
	ctr := controller.NewUserController(ur)

	random := uuid.NewString()
	reqBody := []byte(`{ "email": "existed@test.com"}`)
	req, _ := http.NewRequest("PUT", fmt.Sprintf("/users/%s", random), bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), middleware.RequestIDKey, "test-trace-id")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ctr.UpdateUser)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusConflict, rr.Code)
}

func TestUpdateUserInternalServerError(t *testing.T) {
	ur := &userRepositoryMock{shouldReturn500: true}
	ctr := controller.NewUserController(ur)

	random := uuid.NewString()
	reqBody := []byte(fmt.Sprintf(`{
		"username": "%s",
		"email": "new_%s@test.com",
		"isEnabled": true,
		"password": "%s_super-secret-password"
	  }`, random, random, random))
	req, _ := http.NewRequest("PUT", fmt.Sprintf("/users/%s", random), bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), middleware.RequestIDKey, "test-trace-id")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ctr.UpdateUser)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	var errRes dto.ErrorResponse
	json.Unmarshal(rr.Body.Bytes(), &errRes)
	assert.NotEmpty(t, errRes.Code)
	assert.NotEmpty(t, errRes.Message)
	assert.NotEmpty(t, errRes.TraceId)
}

func TestListUsers(t *testing.T) {
	ur := &userRepositoryMock{}
	ctr := controller.NewUserController(ur)

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
	ur := &userRepositoryMock{totalElements: totalElements}
	ctr := controller.NewUserController(ur)

	size := 2
	page := 0

	req, _ := http.NewRequest("GET", fmt.Sprintf("/users?page=%d&size=%d", page, size), nil)
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ctr.ListUser)
	handler.ServeHTTP(rr, req)

	var userResponseList dto.UserResponseList
	json.Unmarshal(rr.Body.Bytes(), &userResponseList)

	assert.Equal(t, http.StatusOK, rr.Code)

	assert.Equal(t, page, userResponseList.Page)
	assert.Equal(t, size, userResponseList.Size)
	assert.Equal(t, totalElements, userResponseList.TotalElements)
	assert.True(t, strings.HasSuffix(userResponseList.Next, fmt.Sprintf("/users?page=%d&size=%d", page+1, size)))
	assert.Equal(t, size, len(userResponseList.Content))
}

func TestDeleteSuccess(t *testing.T) {
	ur := &userRepositoryMock{}
	ctr := controller.NewUserController(ur)

	random := uuid.NewString()
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/users/%s", random), nil)
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ctr.DeleteUser)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)
}

func TestDeleteInternalServerError(t *testing.T) {
	ur := &userRepositoryMock{shouldReturn500: true}
	ctr := controller.NewUserController(ur)

	random := uuid.NewString()
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/users/%s", random), nil)
	rr := httptest.NewRecorder()
	ctx := context.WithValue(req.Context(), middleware.RequestIDKey, "test-trace-id")
	req = req.WithContext(ctx)

	handler := http.HandlerFunc(ctr.DeleteUser)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	var errRes dto.ErrorResponse
	json.Unmarshal(rr.Body.Bytes(), &errRes)
	assert.NotEmpty(t, errRes.Code)
	assert.NotEmpty(t, errRes.Message)
	assert.NotEmpty(t, errRes.TraceId)
}

type userRepositoryMock struct {
	shouldReturn500 bool
	shouldReturn409 bool
	shouldReturn404 bool
	totalElements   int
}

func (u *userRepositoryMock) Save(user *model.User) error {
	if u.shouldReturn409 {
		return mongo.ErrInvalidIndexValue
	}

	if u.shouldReturn500 {
		return mongo.ErrClientDisconnected
	}

	return nil
}

func (u *userRepositoryMock) FindById(id string) (user *model.User, err error) {
	if u.shouldReturn500 {
		return nil, mongo.ErrClientDisconnected
	}

	if u.shouldReturn404 {
		return nil, mongo.ErrNoDocuments
	}

	return createRandomUser(id, ""), nil
}

func (u *userRepositoryMock) FindByEmail(email string) (user *model.User, err error) {
	if u.shouldReturn404 {
		return nil, mongo.ErrNoDocuments
	}

	return createRandomUser("", email), nil
}

func (u *userRepositoryMock) Delete(id string) error {
	if u.shouldReturn500 {
		return mongo.ErrClientDisconnected
	}

	return nil
}

func (u *userRepositoryMock) Update(user *model.User) error {
	if u.shouldReturn409 {
		return mongo.ErrInvalidIndexValue
	}

	if u.shouldReturn500 {
		return mongo.ErrClientDisconnected
	}

	return nil

}

func (u *userRepositoryMock) ExistsByEmailOrUsername(email string, username string) (bool, error) {

	if u.shouldReturn500 {
		return false, mongo.ErrClientDisconnected
	}

	if u.shouldReturn409 {
		return true, nil
	}

	return false, nil
}

func (u *userRepositoryMock) FindAll(page int, size int) (userPage *model.UserPage, err error) {

	users := make([]*model.User, 0)

	for i := 0; i < size; i++ {
		users = append(users, createRandomUser("", ""))
	}

	userPage = &model.UserPage{
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
