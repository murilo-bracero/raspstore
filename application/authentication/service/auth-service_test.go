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
	fire, errFire := db.NewFirebaseConnection(context.Background())
	assert.NoError(t, err)
	assert.NoError(t, errFire)
	defer conn.Close(context.Background())
	repo := repository.NewMongoUsersRepository(context.Background(), conn)
	cred := repository.NewFireCredentials(context.Background(), fire)
	service := NewAuthService(repo, cred)

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
	fire, errFire := db.NewFirebaseConnection(context.Background())
	assert.NoError(t, err)
	assert.NoError(t, errFire)
	defer conn.Close(context.Background())
	repo := repository.NewMongoUsersRepository(context.Background(), conn)
	cred := repository.NewFireCredentials(context.Background(), fire)
	service := NewAuthService(repo, cred)

	user := &model.User{
		Username:    fmt.Sprintf("tes_%s", uuid.NewString()),
		Email:       fmt.Sprintf("%s@email.com", uuid.NewString()),
		PhoneNumber: "+552738361319",
	}

	err = repo.Save(user)

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
	fire, errFire := db.NewFirebaseConnection(context.Background())
	assert.NoError(t, err)
	assert.NoError(t, errFire)
	defer conn.Close(context.Background())
	repo := repository.NewMongoUsersRepository(context.Background(), conn)
	cred := repository.NewFireCredentials(context.Background(), fire)
	service := NewAuthService(repo, cred)

	user := &model.User{
		Username:    fmt.Sprintf("tes_%s", uuid.NewString()),
		Email:       fmt.Sprintf("%s@email.com", uuid.NewString()),
		PhoneNumber: "+552738361320",
	}

	err = repo.Save(user)

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
	fire, errFire := db.NewFirebaseConnection(context.Background())
	assert.NoError(t, err)
	assert.NoError(t, errFire)
	defer conn.Close(context.Background())
	repo := repository.NewMongoUsersRepository(context.Background(), conn)
	cred := repository.NewFireCredentials(context.Background(), fire)
	service := NewAuthService(repo, cred)

	users, err1 := repo.FindAll()

	assert.NoError(t, err1)

	for _, user := range users {
		req := &pb.GetUserRequest{Id: user.UserId}
		service.DeleteUser(context.Background(), req)
	}
}

// DATASTORE TESTS

func TestDsSignUp(t *testing.T) {
	cfg := db.NewConfig()
	conn, err := db.NewMongoConnection(context.Background(), cfg)
	fire, errFire := db.NewFirebaseConnection(context.Background())
	assert.NoError(t, err)
	assert.NoError(t, errFire)
	defer conn.Close(context.Background())
	repo := repository.NewMongoUsersRepository(context.Background(), conn)
	cred := repository.NewFireCredentials(context.Background(), fire)
	service := NewAuthService(repo, cred)

	req := &pb.CreateUserRequest{
		Password:    "testpass",
		Username:    fmt.Sprintf("tes_%s", uuid.NewString()),
		Email:       fmt.Sprintf("%s@email.com", uuid.NewString()),
		PhoneNumber: "+552738361322",
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
	conn, err := db.NewMongoConnection(context.Background(), cfg)
	fire, errFire := db.NewFirebaseConnection(context.Background())
	assert.NoError(t, err)
	assert.NoError(t, errFire)
	defer conn.Close(context.Background())
	repo := repository.NewMongoUsersRepository(context.Background(), conn)
	cred := repository.NewFireCredentials(context.Background(), fire)
	service := NewAuthService(repo, cred)

	user := &model.User{
		Username:    fmt.Sprintf("tes_%s", uuid.NewString()),
		Email:       fmt.Sprintf("%s@email.com", uuid.NewString()),
		PhoneNumber: "+552738361323",
	}

	err = repo.Save(user)

	assert.NoError(t, err)

	req := &pb.GetUserRequest{Id: user.UserId}

	found, err1 := service.GetUser(context.Background(), req)

	assert.NoError(t, err1)
	assert.Equal(t, user.Username, found.Username)
	assert.Equal(t, user.Email, found.Email)
	assert.Equal(t, user.PhoneNumber, found.PhoneNumber)
}

func TestDsUpdateUser(t *testing.T) {
	cfg := db.NewConfig()
	conn, err := db.NewMongoConnection(context.Background(), cfg)
	fire, errFire := db.NewFirebaseConnection(context.Background())
	assert.NoError(t, err)
	assert.NoError(t, errFire)
	defer conn.Close(context.Background())
	repo := repository.NewMongoUsersRepository(context.Background(), conn)
	cred := repository.NewFireCredentials(context.Background(), fire)
	service := NewAuthService(repo, cred)

	user := &model.User{
		Username:    fmt.Sprintf("tes_%s", uuid.NewString()),
		Email:       fmt.Sprintf("%s@email.com", uuid.NewString()),
		PhoneNumber: "+552738361324",
	}

	err = repo.Save(user)

	assert.NoError(t, err)

	req := &pb.UpdateUserRequest{
		Id:          user.UserId,
		Username:    fmt.Sprintf("updated_%s", uuid.NewString()),
		Email:       fmt.Sprintf("updated_%s@email.com", uuid.NewString()),
		PhoneNumber: "+552738361325",
	}

	found, err1 := service.UpdateUser(context.Background(), req)

	assert.NoError(t, err1)
	assert.NotEqual(t, user.Username, found.Username)
	assert.NotEqual(t, user.Email, found.Email)
	assert.Equal(t, user.PhoneNumber, found.PhoneNumber)
}

func TestDsDeleteUser(t *testing.T) {
	cfg := db.NewConfig()
	conn, err := db.NewMongoConnection(context.Background(), cfg)
	fire, errFire := db.NewFirebaseConnection(context.Background())
	assert.NoError(t, err)
	assert.NoError(t, errFire)
	defer conn.Close(context.Background())
	repo := repository.NewMongoUsersRepository(context.Background(), conn)
	cred := repository.NewFireCredentials(context.Background(), fire)
	service := NewAuthService(repo, cred)

	users, err1 := repo.FindAll()

	assert.NoError(t, err1)

	for _, user := range users {
		req := &pb.GetUserRequest{Id: user.UserId}
		service.DeleteUser(context.Background(), req)
	}
}
