package repository_test

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/pquerna/otp/totp"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
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
		Id:     id,
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
		Id:     id,
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

func TestCredsRepositoryIsCredentialsCorrect(t *testing.T) {
	cfg := db.NewConfig()
	conn, err := db.NewMongoConnection(context.Background(), cfg)
	assert.NoError(t, err)
	defer conn.Close(context.Background())
	repo := mg.NewCredentialsRepository(context.Background(), conn)

	id := uuid.NewString()

	hash, err := hash("testpass")

	assert.NoError(t, err)

	secret, err := totp.Generate(totp.GenerateOpts{Issuer: "rasptore", AccountName: fmt.Sprintf("%s@test.com", id)})

	assert.NoError(t, err)

	cred := &model.Credential{
		Id:     id,
		Email:  fmt.Sprintf("%s@test.com", id),
		Secret: secret.Secret(),
		Hash:   hash,
	}

	err = repo.Save(cred)

	assert.NoError(t, err)

	token, err := totp.GenerateCode(cred.Secret, time.Now())

	assert.NoError(t, err)

	isCorrect := repo.IsCredentialsCorrect(cred.Email, "testpass", token)

	assert.True(t, isCorrect)
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
		if err = repo.Delete(cred.Id); err != nil {
			assert.Fail(t, err.Error())
		}
	}
}

func hash(text string) (hash string, err error) {
	bts := []byte(text)

	raw, err := bcrypt.GenerateFromPassword(bts, bcrypt.DefaultCost)
	return string(raw[:]), err
}
