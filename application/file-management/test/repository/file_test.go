package repository_test

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"raspstore.github.io/file-manager/db"
	"raspstore.github.io/file-manager/model"
	"raspstore.github.io/file-manager/repository"
)

func init() {

	err := godotenv.Load("../../.env")

	if err != nil {
		log.Panicln(err.Error())
	}
}

func TestSave(t *testing.T) {
	cfg := db.NewConfig()
	conn, err := db.NewMongoConnection(context.Background(), cfg)
	assert.NoError(t, err)

	fr := repository.NewFilesRepository(context.Background(), conn)

	file := &model.File{
		Id:        primitive.NewObjectID(),
		Filename:  "test.toml",
		Uri:       "test/test.toml",
		Size:      39201,
		UpdatedAt: time.Now(),
		CreatedBy: "test",
		UpdatedBy: "test",
	}

	fr.Save(file)

	found, err := fr.FindById(file.Id.Hex())

	assert.NoError(t, err)

	assert.Equal(t, found.Filename, file.Filename)
}

func TestUpdate(t *testing.T) {
	cfg := db.NewConfig()
	conn, err := db.NewMongoConnection(context.Background(), cfg)
	assert.NoError(t, err)

	fr := repository.NewFilesRepository(context.Background(), conn)

	file := &model.File{
		Id:        primitive.NewObjectID(),
		Filename:  "test.toml",
		Uri:       "test/test.toml",
		Size:      39201,
		UpdatedAt: time.Now(),
		CreatedBy: "test",
		UpdatedBy: "test",
	}

	fr.Save(file)

	file.UpdatedBy = "test1"
	file.Size = 823091809

	err = fr.Update(file)

	assert.NoError(t, err)
}
