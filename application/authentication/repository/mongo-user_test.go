package repository

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"raspstore.github.io/authentication/db"
	"raspstore.github.io/authentication/model"
)

func init() {
	err := godotenv.Load("../.env")

	if err != nil {
		log.Panicln(err.Error())
	}
}

func TestUsersRepositorySave(t *testing.T) {
	cfg := db.NewConfig()
	conn, err := db.NewMongoConnection(context.Background(), cfg)
	assert.NoError(t, err)
	defer conn.Close(context.Background())
	repo := NewMongoUsersRepository(context.Background(), conn)

	id := primitive.NewObjectID()

	user := &model.User{
		Id:          id,
		Email:       fmt.Sprintf("%s@email.com", id.Hex()),
		Username:    "testing_test",
		PhoneNumber: "1196726372912",
	}

	repo.Save(user)

	assert.Equal(t, time.Now().Hour(), user.CreatedAt.Hour())
	assert.Equal(t, time.Now().Hour(), user.UpdatedAt.Hour())
	assert.NotNil(t, user.Id)
}

func TestUsersRepositoryFindById(t *testing.T) {
	cfg := db.NewConfig()
	conn, err := db.NewMongoConnection(context.Background(), cfg)
	assert.NoError(t, err)
	defer conn.Close(context.Background())
	repo := NewMongoUsersRepository(context.Background(), conn)

	user := &model.User{
		Email:       "random@email.com",
		Username:    "testing_test",
		PhoneNumber: "1196726372912",
	}

	err = repo.Save(user)

	assert.NoError(t, err)

	found, err1 := repo.FindById(user.Id.Hex())
	assert.NoError(t, err1)
	assert.NotNil(t, found)
}

func TestUsersRepositoryFindByEmailOrUsername(t *testing.T) {
	cfg := db.NewConfig()
	conn, err := db.NewMongoConnection(context.Background(), cfg)
	assert.NoError(t, err)
	defer conn.Close(context.Background())
	repo := NewMongoUsersRepository(context.Background(), conn)

	id := primitive.NewObjectID()

	email := fmt.Sprintf("%s@email.com", id.Hex())
	username := "testing_test"

	user := &model.User{
		Id:          id,
		Email:       email,
		Username:    username,
		PhoneNumber: "1196726372912",
	}

	repo.Save(user)

	found, err1 := repo.FindByEmailOrUsername(email, username)
	assert.NoError(t, err1)
	assert.NotNil(t, found)
}

func TestUsersRepositoryUpdateUser(t *testing.T) {
	cfg := db.NewConfig()
	conn, err := db.NewMongoConnection(context.Background(), cfg)
	assert.NoError(t, err)
	defer conn.Close(context.Background())
	repo := NewMongoUsersRepository(context.Background(), conn)

	user := &model.User{
		Email:       "test@email.com",
		Username:    "testing_test",
		PhoneNumber: "1196726372912",
	}

	repo.Save(user)

	updated_email := fmt.Sprintf("updated_%s@email.com", user.Id.Hex())

	updated := &model.User{
		Id:          user.Id,
		Email:       updated_email,
		Username:    "testing_test",
		PhoneNumber: "1196726372912",
	}

	error1 := repo.UpdateUser(updated)

	assert.NoError(t, error1)

	found, error2 := repo.FindById(user.Id.Hex())

	assert.NoError(t, error2)

	assert.Equal(t, updated_email, found.Email)
}

func TestUsersRepositoryFindAll(t *testing.T) {
	cfg := db.NewConfig()
	conn, err := db.NewMongoConnection(context.Background(), cfg)
	assert.NoError(t, err)
	defer conn.Close(context.Background())
	repo := NewMongoUsersRepository(context.Background(), conn)

	users, error1 := repo.FindAll()
	assert.NoError(t, error1)
	assert.True(t, len(users) > 0)
}

func TestUsersRepositoryDeleteUser(t *testing.T) {
	cfg := db.NewConfig()
	conn, err := db.NewMongoConnection(context.Background(), cfg)
	assert.NoError(t, err)
	defer conn.Close(context.Background())
	repo := NewMongoUsersRepository(context.Background(), conn)

	users, error1 := repo.FindAll()

	assert.NoError(t, error1)
	assert.True(t, len(users) > 0)

	for _, user := range users {
		repo.DeleteUser(user.Id.Hex())
	}
}
