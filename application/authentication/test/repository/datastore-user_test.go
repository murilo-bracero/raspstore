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
	ds "raspstore.github.io/authentication/repository/datastore"
)

func init() {
	err := godotenv.Load("../.env")

	if err != nil {
		log.Panicln(err.Error())
	}
}

func TestDsUsersRepositorySave(t *testing.T) {
	cfg := db.NewConfig()
	conn, err := db.NewDatastoreConnection(context.Background(), cfg)
	assert.NoError(t, err)
	defer conn.Close()
	repo := ds.NewDatastoreUsersRepository(context.Background(), conn)

	id := primitive.NewObjectID()

	user := &model.User{
		Email:       fmt.Sprintf("%s@email.com", id.Hex()),
		Username:    "testing_test",
		PhoneNumber: "1196726372912",
	}

	err = repo.Save(user)

	assert.NoError(t, err)
	assert.Equal(t, time.Now().Hour(), user.CreatedAt.Hour())
	assert.Equal(t, time.Now().Hour(), user.UpdatedAt.Hour())
	assert.NotNil(t, user.UserId)
}

func TestDsUsersRepositoryFindById(t *testing.T) {
	cfg := db.NewConfig()
	conn, err := db.NewDatastoreConnection(context.Background(), cfg)
	assert.NoError(t, err)
	defer conn.Close()
	repo := ds.NewDatastoreUsersRepository(context.Background(), conn)

	id := primitive.NewObjectID()

	user := &model.User{
		Email:       fmt.Sprintf("%s@email.com", id.Hex()),
		Username:    "testing_test",
		PhoneNumber: "1196726372912",
	}

	err = repo.Save(user)

	assert.NoError(t, err)

	found, err1 := repo.FindById(user.UserId)
	assert.NoError(t, err1)
	assert.NotNil(t, found)
}

func TestDsUsersRepositoryFindByEmailOrUsername(t *testing.T) {
	cfg := db.NewConfig()
	conn, err := db.NewDatastoreConnection(context.Background(), cfg)
	assert.NoError(t, err)
	defer conn.Close()
	repo := ds.NewDatastoreUsersRepository(context.Background(), conn)

	email := "TestDsUsersRepositoryFindByEmailOrUsername@email.com"
	username := "testing_test"

	user := &model.User{
		UserId:      "b8fa44f8425640d4a3554518ad5a97ab",
		Email:       email,
		Username:    username,
		PhoneNumber: "1196726372912",
	}

	repo.Save(user)

	found, err1 := repo.FindByEmail(email)
	assert.NoError(t, err1)
	assert.NotNil(t, found)
}

func TestDsUsersRepositoryUpdateUser(t *testing.T) {
	cfg := db.NewConfig()
	conn, err := db.NewDatastoreConnection(context.Background(), cfg)
	assert.NoError(t, err)
	defer conn.Close()
	repo := ds.NewDatastoreUsersRepository(context.Background(), conn)

	id := primitive.NewObjectID()

	user := &model.User{
		Email:       fmt.Sprintf("%s@email.com", id.Hex()),
		Username:    "testing_test",
		PhoneNumber: "1196726372912",
	}

	err = repo.Save(user)

	assert.NoError(t, err)

	updated_email := fmt.Sprintf("updated_%s@email.com", id.Hex())

	updated := &model.User{
		UserId:      user.UserId,
		Email:       updated_email,
		Username:    "testing_test",
		PhoneNumber: "1196726372912",
	}

	err = repo.UpdateUser(updated)

	assert.NoError(t, err)

	found, err := repo.FindById(user.UserId)

	assert.NoError(t, err)

	assert.Equal(t, updated_email, found.Email)
}

func TestDsUsersRepositoryFindAll(t *testing.T) {
	cfg := db.NewConfig()
	conn, err := db.NewDatastoreConnection(context.Background(), cfg)
	assert.NoError(t, err)
	defer conn.Close()
	repo := ds.NewDatastoreUsersRepository(context.Background(), conn)

	users, error1 := repo.FindAll()
	assert.NoError(t, error1)
	assert.True(t, len(users) > 0)
}

func TestDsUsersRepositoryDeleteUser(t *testing.T) {
	cfg := db.NewConfig()
	conn, err := db.NewDatastoreConnection(context.Background(), cfg)
	assert.NoError(t, err)
	defer conn.Close()
	repo := ds.NewDatastoreUsersRepository(context.Background(), conn)

	users, error1 := repo.FindAll()

	assert.NoError(t, error1)
	assert.True(t, len(users) > 0)

	for _, user := range users {
		repo.DeleteUser(user.UserId)
	}
}
