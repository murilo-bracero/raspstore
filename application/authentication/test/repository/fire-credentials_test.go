package repository

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"google.golang.org/api/iterator"
	"raspstore.github.io/authentication/db"
	"raspstore.github.io/authentication/model"
	fb "raspstore.github.io/authentication/repository/firebase"
)

func init() {
	err := godotenv.Load("../.env")

	if err != nil {
		log.Panicln(err.Error())
	}
}

func TestFireCredentialsRepositorySave(t *testing.T) {
	conn, err := db.NewFirebaseConnection(context.Background())
	assert.NoError(t, err)
	repo := fb.NewFireCredentials(context.Background(), conn)

	user := &model.User{
		UserId:      uuid.NewString(),
		Username:    fmt.Sprintf("test_%s", uuid.NewString()),
		Email:       fmt.Sprintf("test_%s@email.com", uuid.NewString()),
		PhoneNumber: "+554367678989",
	}

	err = repo.Save(user, "test14235")
	assert.NoError(t, err)
}

func TestFireCredentialsRepositoryUpdate(t *testing.T) {
	conn, err := db.NewFirebaseConnection(context.Background())
	assert.NoError(t, err)
	repo := fb.NewFireCredentials(context.Background(), conn)

	user := &model.User{
		UserId:      uuid.NewString(),
		Username:    fmt.Sprintf("test_%s", uuid.NewString()),
		Email:       fmt.Sprintf("test_%s@email.com", uuid.NewString()),
		PhoneNumber: "+554367678990",
	}

	err = repo.Save(user, "test14235")
	assert.NoError(t, err)

	user.Username = fmt.Sprintf("updated_%s", uuid.NewString())

	err = repo.Update(user)
	assert.NoError(t, err)
}

func TestFireCredentialsRepositoryDelete(t *testing.T) {
	conn, err := db.NewFirebaseConnection(context.Background())
	assert.NoError(t, err)
	repo := fb.NewFireCredentials(context.Background(), conn)

	iter := conn.Client().Users(context.Background(), "")

	for {
		user, errIter := iter.Next()

		if errIter == iterator.Done {
			break
		}

		assert.NoError(t, errIter)
		err = repo.Delete(user.UID)
		assert.NoError(t, err)
	}
}
