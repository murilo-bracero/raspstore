package service_test

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"raspstore.github.io/users-service/db"
	"raspstore.github.io/users-service/model"
	"raspstore.github.io/users-service/pb"
	mg "raspstore.github.io/users-service/repository"
	sv "raspstore.github.io/users-service/service"
)

func init() {
	err := godotenv.Load("../../.env")

	if err != nil {
		log.Panicln(err.Error())
	}
}

func TestSignUp(t *testing.T) {
	cfg := db.NewConfig()
	conn, err := db.NewMongoConnection(context.Background(), cfg)
	assert.NoError(t, err)
	defer conn.Close(context.Background())
	userRepo := mg.NewUsersRepository(context.Background(), conn)
	credRepo := mg.NewCredentialsRepository(context.Background(), conn)
	service := sv.NewUserService(userRepo, credRepo)

	req := &pb.CreateUserRequest{
		Username:    fmt.Sprintf("tes_%s", uuid.NewString()),
		Password:    "testpass",
		Email:       fmt.Sprintf("%s@email.com", uuid.NewString()),
		PhoneNumber: "+552738361318",
	}

	user, err := service.CreateUser(context.Background(), req)

	assert.NoError(t, err)
	assert.Equal(t, req.Username, user.Username)
	assert.Equal(t, req.Email, user.Email)
	assert.Equal(t, req.PhoneNumber, user.PhoneNumber)
	assert.NotNil(t, user.CreatedAt)
	assert.NotNil(t, user.UpdatedAt)
}

func TestGetUser(t *testing.T) {
	cfg := db.NewConfig()
	conn, err := db.NewMongoConnection(context.Background(), cfg)
	assert.NoError(t, err)
	defer conn.Close(context.Background())
	userRepo := mg.NewUsersRepository(context.Background(), conn)
	credRepo := mg.NewCredentialsRepository(context.Background(), conn)
	service := sv.NewUserService(userRepo, credRepo)

	user := &model.User{
		Username:    fmt.Sprintf("tes_%s", uuid.NewString()),
		Email:       fmt.Sprintf("%s@email.com", uuid.NewString()),
		PhoneNumber: "+552738361319",
	}

	err = userRepo.Save(user)

	assert.NoError(t, err)

	req := &pb.GetUserRequest{Id: user.UserId}

	found, err1 := service.GetUser(context.Background(), req)

	assert.NoError(t, err1)
	assert.Equal(t, user.Username, found.Username)
	assert.Equal(t, user.Email, found.Email)
	assert.Equal(t, user.PhoneNumber, found.PhoneNumber)
}

func TestUpdateUser(t *testing.T) {
	cfg := db.NewConfig()
	conn, err := db.NewMongoConnection(context.Background(), cfg)
	assert.NoError(t, err)
	defer conn.Close(context.Background())
	userRepo := mg.NewUsersRepository(context.Background(), conn)
	credRepo := mg.NewCredentialsRepository(context.Background(), conn)
	service := sv.NewUserService(userRepo, credRepo)

	createUserRequest := &pb.CreateUserRequest{
		Username:    fmt.Sprintf("tes_%s", uuid.NewString()),
		Email:       fmt.Sprintf("%s@email.com", uuid.NewString()),
		PhoneNumber: "+552738361320",
		Password:    "penispintorola212",
	}

	user, errService := service.CreateUser(context.Background(), createUserRequest)

	assert.NoError(t, errService)

	req := &pb.UpdateUserRequest{
		Id:          user.Id,
		Username:    fmt.Sprintf("updated_%s", uuid.NewString()),
		Email:       fmt.Sprintf("updated_%s@email.com", uuid.NewString()),
		PhoneNumber: "+552738361321",
	}

	found, err1 := service.UpdateUser(context.Background(), req)

	assert.NoError(t, err1)
	assert.NotEqual(t, user.Username, found.Username)
	assert.NotEqual(t, user.Email, found.Email)
	assert.NotEqual(t, user.PhoneNumber, found.PhoneNumber)
}

func TestDeleteUser(t *testing.T) {
	cfg := db.NewConfig()
	conn, err := db.NewMongoConnection(context.Background(), cfg)
	assert.NoError(t, err)
	defer conn.Close(context.Background())
	userRepo := mg.NewUsersRepository(context.Background(), conn)
	credRepo := mg.NewCredentialsRepository(context.Background(), conn)
	service := sv.NewUserService(userRepo, credRepo)

	users, err1 := userRepo.FindAll()

	assert.NoError(t, err1)

	for _, user := range users {
		req := &pb.GetUserRequest{Id: user.UserId}
		service.DeleteUser(context.Background(), req)
	}
}
