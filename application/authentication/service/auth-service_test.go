package service

import (
	"context"
	"log"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"raspstore.github.io/authentication/db"
	"raspstore.github.io/authentication/model"
	"raspstore.github.io/authentication/pb"
	"raspstore.github.io/authentication/repository"
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
	repo := repository.NewMongoUsersRepository(context.Background(), conn)
	service := NewAuthService(repo)

	req := &pb.CreateUserRequest{
		Username:    "test",
		Password:    "testpass",
		Email:       "test@email.com",
		PhoneNumber: "2783918273",
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
	repo := repository.NewMongoUsersRepository(context.Background(), conn)

	user := &model.User{
		Username:    "testinghehe",
		Email:       "testing@test.com.test",
		PhoneNumber: "39820129021",
	}

	err = repo.Save(user)

	assert.NoError(t, err)

	service := NewAuthService(repo)

	req := &pb.GetUserRequest{Id: user.Id.Hex()}

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
	repo := repository.NewMongoUsersRepository(context.Background(), conn)

	user := &model.User{
		Username:    "testinghehe",
		Email:       "testing@test.com.test",
		PhoneNumber: "39820129021",
	}

	err = repo.Save(user)

	assert.NoError(t, err)

	service := NewAuthService(repo)

	req := &pb.UpdateUserRequest{
		Id:          user.Id.Hex(),
		Username:    "updated_spookyscary",
		Email:       "updated_spookyscary@email.com",
		PhoneNumber: "39820129021",
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
	repo := repository.NewMongoUsersRepository(context.Background(), conn)

	users, err1 := repo.FindAll()

	assert.NoError(t, err1)

	service := NewAuthService(repo)

	for _, user := range users {
		req := &pb.GetUserRequest{Id: user.Id.Hex()}
		service.DeleteUser(context.Background(), req)
	}
}

// DATASTORE TESTS

func TestDsSignUp(t *testing.T) {
	cfg := db.NewConfig()
	conn, err := db.NewDatastoreConnection(context.Background(), cfg)
	assert.NoError(t, err)
	defer conn.Close()
	repo := repository.NewDatastoreUsersRepository(context.Background(), conn)
	service := NewAuthService(repo)

	req := &pb.CreateUserRequest{
		Username:    "test",
		Password:    "testpass",
		Email:       "test@email.com",
		PhoneNumber: "2783918273",
	}

	user, err := service.SignUp(context.Background(), req)

	assert.NoError(t, err)
	assert.Equal(t, req.Username, user.Username)
	assert.Equal(t, req.Email, user.Email)
	assert.Equal(t, req.PhoneNumber, user.PhoneNumber)
	assert.NotNil(t, user.CreatedAt)
	assert.NotNil(t, user.UpdatedAt)
}

func TestDsGetUser(t *testing.T) {
	cfg := db.NewConfig()
	conn, err := db.NewDatastoreConnection(context.Background(), cfg)
	assert.NoError(t, err)
	defer conn.Close()
	repo := repository.NewDatastoreUsersRepository(context.Background(), conn)

	user := &model.User{
		Username:    "testinghehe",
		Email:       "testing@test.com.test",
		PhoneNumber: "39820129021",
	}

	err = repo.Save(user)

	assert.NoError(t, err)

	service := NewAuthService(repo)

	req := &pb.GetUserRequest{Id: user.UserId}

	found, err1 := service.GetUser(context.Background(), req)

	assert.NoError(t, err1)
	assert.Equal(t, user.Username, found.Username)
	assert.Equal(t, user.Email, found.Email)
	assert.Equal(t, user.PhoneNumber, found.PhoneNumber)
}

func TestDsUpdateUser(t *testing.T) {
	cfg := db.NewConfig()
	conn, err := db.NewDatastoreConnection(context.Background(), cfg)
	assert.NoError(t, err)
	defer conn.Close()
	repo := repository.NewDatastoreUsersRepository(context.Background(), conn)

	user := &model.User{
		Username:    "testinghehe",
		Email:       "testing@test.com.test",
		PhoneNumber: "39820129021",
	}

	err = repo.Save(user)

	assert.NoError(t, err)

	service := NewAuthService(repo)

	req := &pb.UpdateUserRequest{
		Id:          user.UserId,
		Username:    "updated_spookyscary",
		Email:       "updated_spookyscary@email.com",
		PhoneNumber: "39820129021",
	}

	found, err1 := service.UpdateUser(context.Background(), req)

	assert.NoError(t, err1)
	assert.NotEqual(t, user.Username, found.Username)
	assert.NotEqual(t, user.Email, found.Email)
	assert.Equal(t, user.PhoneNumber, found.PhoneNumber)
}

func TestDsDeleteUser(t *testing.T) {
	cfg := db.NewConfig()
	conn, err := db.NewDatastoreConnection(context.Background(), cfg)
	assert.NoError(t, err)
	defer conn.Close()
	repo := repository.NewDatastoreUsersRepository(context.Background(), conn)

	users, err1 := repo.FindAll()

	assert.NoError(t, err1)

	service := NewAuthService(repo)

	for _, user := range users {
		req := &pb.GetUserRequest{Id: user.UserId}
		service.DeleteUser(context.Background(), req)
	}
}
