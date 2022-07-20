package controller_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"raspstore.github.io/users-service/api"
	"raspstore.github.io/users-service/api/controller"
	"raspstore.github.io/users-service/db"
	"raspstore.github.io/users-service/model"
	"raspstore.github.io/users-service/repository"
	"raspstore.github.io/users-service/service"
)

var routes api.Routes

func TestMain(m *testing.M) {
	err := godotenv.Load("../../../.env")

	if err != nil {
		log.Panicln(err.Error())
	}

	cfg := db.NewConfig()
	conn, err := db.NewMongoConnection(context.Background(), cfg)

	if err != nil {
		log.Panicln(err.Error())
	}

	defer conn.Close(context.Background())
	ur := repository.NewUsersRepository(context.Background(), conn)
	svc := service.NewUserService(ur)
	uc := controller.NewUserController(ur, svc)
	routes = api.NewRoutes(uc)
	code := m.Run()
	err = teardown(ur)
	if err != nil {
		log.Panicln(err.Error())
	}
	os.Exit(code)
}

func TestCreateUserSuccess(t *testing.T) {
	random := uuid.NewString()
	random = strings.ReplaceAll(random, "-", "")

	raw := fmt.Sprintf(`{
		"username": "%s-username",
		"password":"defaultPassword",
		"email": "%s@fake.com.br",
		"phoneNumber": "+5511243516237"
	}`, random, random)

	body := []byte(raw)

	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rawRes := sendReq(req)

	assert.Equal(t, rawRes.Code, 201)

	var res model.User
	json.Unmarshal(rawRes.Body.Bytes(), &res)

	assert.NotNil(t, res.CreatedAt)
	assert.NotNil(t, res.UpdatedAt)
	assert.Equal(t, res.CreatedAt, res.UpdatedAt)
	assert.NotNil(t, res.UserId)
	assert.Equal(t, fmt.Sprintf("%s-username", random), res.Username)
	assert.Equal(t, fmt.Sprintf("%s@fake.com.br", random), res.Email)
}

func TestCreateUserFailOnAlreadyExistedEmail(t *testing.T) {
	random := uuid.NewString()
	random = strings.ReplaceAll(random, "-", "")

	raw := fmt.Sprintf(`{
		"username": "%s-username",
		"password":"defaultPassword",
		"email": "%s@fake.com.br",
		"phoneNumber": "+5511243516237"
	}`, random, random)

	body := []byte(raw)

	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rawRes := sendReq(req)

	assert.Equal(t, rawRes.Code, 201)

	req, _ = http.NewRequest("POST", "/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rawRes = sendReq(req)

	assert.Equal(t, rawRes.Code, 400)
}

func TestCreateUserFailOnIncompletePayload(t *testing.T) {
	random := uuid.NewString()
	random = strings.ReplaceAll(random, "-", "")

	raw := fmt.Sprintf(`{
		"username": "%s-username",
		"password":"defaultPassword",
		"phoneNumber": "+5511243516237"
	}`, random)

	body := []byte(raw)

	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rawRes := sendReq(req)

	assert.Equal(t, rawRes.Code, 400)
}

func TestGetUser(t *testing.T) {
	random := uuid.NewString()
	random = strings.ReplaceAll(random, "-", "")

	raw := fmt.Sprintf(`{
		"username": "%s-username",
		"password":"defaultPassword",
		"email": "%s@fake.com.br",
		"phoneNumber": "+5511243516237"
	}`, random, random)

	body := []byte(raw)

	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rawRes := sendReq(req)

	assert.Equal(t, rawRes.Code, 201)

	var res model.User
	json.Unmarshal(rawRes.Body.Bytes(), &res)

	assert.NotNil(t, res.UserId)

	req, _ = http.NewRequest("GET", fmt.Sprintf("/users/%s", res.UserId), nil)
	req.Header.Set("Content-Type", "application/json")

	rawRes = sendReq(req)

	assert.Equal(t, rawRes.Code, 200)

	json.Unmarshal(rawRes.Body.Bytes(), &res)

	assert.NotNil(t, res.CreatedAt)
	assert.NotNil(t, res.UpdatedAt)
	assert.Equal(t, res.CreatedAt, res.UpdatedAt)
	assert.NotNil(t, res.UserId)
	assert.Equal(t, fmt.Sprintf("%s-username", random), res.Username)
	assert.Equal(t, fmt.Sprintf("%s@fake.com.br", random), res.Email)
}

func TestGetUserWithInexistedId(t *testing.T) {
	random := uuid.NewString()

	req, _ := http.NewRequest("GET", fmt.Sprintf("/users/%s", random), nil)
	req.Header.Set("Content-Type", "application/json")

	rawRes := sendReq(req)

	assert.Equal(t, rawRes.Code, 404)
}

func TestUpdateUserSuccess(t *testing.T) {
	random := uuid.NewString()
	random = strings.ReplaceAll(random, "-", "")

	raw := fmt.Sprintf(`{
		"username": "%s-username",
		"password":"defaultPassword",
		"email": "%s@fake.com.br",
		"phoneNumber": "+5511243516237"
	}`, random, random)

	body := []byte(raw)

	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rawRes := sendReq(req)

	assert.Equal(t, rawRes.Code, 201)

	var res model.User
	json.Unmarshal(rawRes.Body.Bytes(), &res)

	assert.NotNil(t, res.UserId)

	raw = fmt.Sprintf(`{
		"username": "updated_%s-username",
		"email": "updated_%s@fake.com.br",
		"phoneNumber": "+5511243516238"
	}`, random, random)

	body = []byte(raw)

	req, _ = http.NewRequest("PATCH", fmt.Sprintf("/users/%s", res.UserId), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rawRes = sendReq(req)

	assert.Equal(t, rawRes.Code, 200)

	json.Unmarshal(rawRes.Body.Bytes(), &res)

	assert.True(t, strings.Contains(res.Username, "updated"))
	assert.True(t, strings.Contains(res.Email, "updated"))

}

func TestUpdateUserFailOnIncompletePayload(t *testing.T) {
	random := uuid.NewString()
	random = strings.ReplaceAll(random, "-", "")

	raw := fmt.Sprintf(`{
		"username": "%s-username",
		"password":"defaultPassword",
		"email": "%s@fake.com.br",
		"phoneNumber": "+5511243516237"
	}`, random, random)

	body := []byte(raw)

	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rawRes := sendReq(req)

	assert.Equal(t, rawRes.Code, 201)

	var res model.User
	json.Unmarshal(rawRes.Body.Bytes(), &res)

	assert.NotNil(t, res.UserId)

	raw = fmt.Sprintf(`{
		"username": "updated_%s-username",
		"phoneNumber": "+5511243516238"
	}`, random)

	body = []byte(raw)

	req, _ = http.NewRequest("PATCH", fmt.Sprintf("/users/%s", res.UserId), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rawRes = sendReq(req)

	assert.Equal(t, rawRes.Code, 400)
}

func TestUpdateUserFailOnInexistentUserId(t *testing.T) {
	random := uuid.NewString()

	raw := fmt.Sprintf(`{
		"username": "%s-username",
		"password":"defaultPassword",
		"email": "%s@fake.com.br",
		"phoneNumber": "+5511243516237"
	}`, random, random)

	body := []byte(raw)

	req, _ := http.NewRequest("PATCH", fmt.Sprintf("/users/%s", random), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rawRes := sendReq(req)

	assert.Equal(t, 404, rawRes.Code)
}

func TestDeleteUser(t *testing.T) {

	random := uuid.NewString()
	random = strings.ReplaceAll(random, "-", "")

	raw := fmt.Sprintf(`{
		"username": "%s-username",
		"password":"defaultPassword",
		"email": "%s@fake.com.br",
		"phoneNumber": "+5511243516237"
	}`, random, random)

	body := []byte(raw)

	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rawRes := sendReq(req)

	assert.Equal(t, rawRes.Code, 201)

	var res model.User
	json.Unmarshal(rawRes.Body.Bytes(), &res)

	assert.NotNil(t, res.UserId)

	req, _ = http.NewRequest("DELETE", fmt.Sprintf("/users/%s", res.UserId), nil)
	req.Header.Set("Content-Type", "application/json")

	rawRes = sendReq(req)

	assert.Equal(t, rawRes.Code, 204)

	req, _ = http.NewRequest("GET", fmt.Sprintf("/users/%s", res.UserId), nil)
	req.Header.Set("Content-Type", "application/json")

	rawRes = sendReq(req)

	assert.Equal(t, 404, rawRes.Code)
}

func TestListUsers(t *testing.T) {

	req, _ := http.NewRequest("GET", "/users", nil)
	res := sendReq(req)

	var users []model.User
	json.Unmarshal(res.Body.Bytes(), &users)

	assert.Equal(t, 200, res.Code)

	assert.True(t, len(users) > 0)
}

func teardown(ur repository.UsersRepository) error {

	users, err := ur.FindAll()

	if err != nil {
		return err
	}

	for _, user := range users {
		ur.DeleteUser(user.UserId)
	}

	return nil
}

func sendReq(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	routes.MountRoutes().ServeHTTP(rr, req)

	return rr
}
