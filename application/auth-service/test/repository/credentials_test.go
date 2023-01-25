package repository_test

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"raspstore.github.io/authentication/db"
	"raspstore.github.io/authentication/model"
	mg "raspstore.github.io/authentication/repository"
)

func init() {
	err := godotenv.Load("../../.env")

	if err != nil {
		log.Panicln(err.Error())
	}
}

func TestCredsRepositorySave(t *testing.T) {
	cfg := db.NewConfig()
	conn, err := db.NewMongoConnection(context.Background(), cfg)
	assert.NoError(t, err)
	defer conn.Close(context.Background())
	repo := mg.NewCredentialsRepository(context.Background(), conn)

	id := uuid.NewString()

	cred := &model.Credential{
		Id:     primitive.NewObjectID(),
		Email:  fmt.Sprintf("%s@test.com", id),
		Secret: "testing_test",
		Hash:   "testing",
	}

	err = repo.Save(cred)

	assert.NoError(t, err)
	assert.NotNil(t, cred.Id)
}

func TestCredsRepositoryUpdate(t *testing.T) {
	cfg := db.NewConfig()
	conn, err := db.NewMongoConnection(context.Background(), cfg)
	assert.NoError(t, err)
	defer conn.Close(context.Background())
	repo := mg.NewCredentialsRepository(context.Background(), conn)

	id := uuid.NewString()

	cred := &model.Credential{
		Id:     primitive.NewObjectID(),
		UserId: "test_user_id",
		Email:  fmt.Sprintf("%s@test.com", id),
		Secret: "testing_test",
		Hash:   "testing",
	}

	err = repo.Save(cred)

	assert.NoError(t, err)

	cred.Email = "sus@sus.com"

	err = repo.Update(cred)

	assert.NoError(t, err)
	assert.NotNil(t, cred.Id)
}

func TestCredsRepositoryDelete(t *testing.T) {
	ctx := context.Background()
	cfg := db.NewConfig()
	conn, err := db.NewMongoConnection(ctx, cfg)
	assert.NoError(t, err)
	defer conn.Close(ctx)
	repo := mg.NewCredentialsRepository(ctx, conn)

	cursor, errFind := conn.DB().Collection("credentials").Find(ctx, bson.M{"email": bson.M{"$regex": "@test.com"}})

	assert.NoError(t, errFind)

	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var cred model.Credential
		if err = cursor.Decode(&cred); err != nil {
			assert.Fail(t, err.Error())
		}
		if err = repo.Delete(cred.Id.Hex()); err != nil {
			assert.Fail(t, err.Error())
		}
	}
}
