package service

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"raspstore.github.io/authentication/db"
	"raspstore.github.io/authentication/model"
	"raspstore.github.io/authentication/pb"
	mg "raspstore.github.io/authentication/repository"
	sv "raspstore.github.io/authentication/service"
)

func init() {
	err := godotenv.Load("../.env")

	if err != nil {
		log.Panicln(err.Error())
	}
}

// MONGO TESTS

func TestMongoSignUp(t *testing.T) {
	cfg := db.NewConfig()
	conn, err := db.NewMongoConnection(context.Background(), cfg)
	assert.NoError(t, err)
	defer conn.Close(context.Background())
	userRepo := mg.NewMongoUsersRepository(context.Background(), conn)
	credRepo := mg.NewMongoCredentialsRepository(context.Background(), conn)
	service := sv.NewAuthService(userRepo, credRepo)

	req := &pb.CreateUserRequest{
		Username:    fmt.Sprintf("tes_%s", uuid.NewString()),
		Password:    "testpass",
		Email:       fmt.Sprintf("%s@email.com", uuid.NewString()),
		PhoneNumber: "+552738361318",
	}

	user, err := service.SignUp(context.Background(), req)

	assert.NoError(t, err)
	assert.Equal(t, req.Username, user.Username)
	assert.Equal(t, req.Email, user.Email)
	assert.Equal(t, req.PhoneNumber, user.PhoneNumber)
	assert.NotNil(t, user.CreatedAt)
	assert.NotNil(t, user.UpdatedAt)
}

func TestMongoGetUser(t *testing.T) {
	cfg := db.NewConfig()
	conn, err := db.NewMongoConnection(context.Background(), cfg)
	assert.NoError(t, err)
	defer conn.Close(context.Background())
	userRepo := mg.NewMongoUsersRepository(context.Background(), conn)
	credRepo := mg.NewMongoCredentialsRepository(context.Background(), conn)
	service := sv.NewAuthService(userRepo, credRepo)

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

func TestMongoUpdateUser(t *testing.T) {
	cfg := db.NewConfig()
	conn, err := db.NewMongoConnection(context.Background(), cfg)
	assert.NoError(t, err)
	defer conn.Close(context.Background())
	userRepo := mg.NewMongoUsersRepository(context.Background(), conn)
	credRepo := mg.NewMongoCredentialsRepository(context.Background(), conn)
	service := sv.NewAuthService(userRepo, credRepo)

	user := &model.User{
		Username:    fmt.Sprintf("tes_%s", uuid.NewString()),
		Email:       fmt.Sprintf("%s@email.com", uuid.NewString()),
		PhoneNumber: "+552738361320",
	}

	err = userRepo.Save(user)

	assert.NoError(t, err)

	req := &pb.UpdateUserRequest{
		Id:          user.UserId,
		Username:    fmt.Sprintf("updated_%s", uuid.NewString()),
		Email:       fmt.Sprintf("updated_%s@email.com", uuid.NewString()),
		PhoneNumber: "+552738361321",
	}

	found, err1 := service.UpdateUser(context.Background(), req)

	assert.NoError(t, err1)
	assert.NotEqual(t, user.Username, found.Username)
	assert.NotEqual(t, user.Email, found.Email)
	assert.Equal(t, user.PhoneNumber, found.PhoneNumber)
}

func TestMongoDeleteUser(t *testing.T) {
	cfg := db.NewConfig()
	conn, err := db.NewMongoConnection(context.Background(), cfg)
	assert.NoError(t, err)
	defer conn.Close(context.Background())
	userRepo := mg.NewMongoUsersRepository(context.Background(), conn)
	credRepo := mg.NewMongoCredentialsRepository(context.Background(), conn)
	service := sv.NewAuthService(userRepo, credRepo)

	users, err1 := userRepo.FindAll()

	assert.NoError(t, err1)

	for _, user := range users {
		req := &pb.GetUserRequest{Id: user.UserId}
		service.DeleteUser(context.Background(), req)
	}
}
