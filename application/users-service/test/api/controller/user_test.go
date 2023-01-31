package controller_test

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
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

	usr := &model.User{
		UserId:      uuid.NewString(),
		Username:    "test username",
		Email:       "test_email@email.com",
		PhoneNumber: "019202012131",
	}

	ur.Save(usr)

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

func TestGetUserWithInexistedId(t *testing.T) {
	random := uuid.NewString()

	req, _ := http.NewRequest("GET", fmt.Sprintf("/users/%s", random), nil)
	req.Header.Set("Content-Type", "application/json")

	rawRes := sendReq(req)

	assert.Equal(t, rawRes.Code, 404)
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
		ur.Delete(user.UserId)
	}

	return nil
}

func sendReq(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	routes.MountRoutes().ServeHTTP(rr, req)

	return rr
}
