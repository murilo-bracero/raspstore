package system_test

import (
	"bytes"
	"fmt"
	"log"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"raspstore.github.io/file-manager/db"
	"raspstore.github.io/file-manager/model"
	"raspstore.github.io/file-manager/system"
)

func init() {

	err := godotenv.Load("../../.env")

	if err != nil {
		log.Panicln(err.Error())
	}
}

func TestSaveFile(t *testing.T) {
	cfg := db.NewConfig()
	ds := system.NewDiskStore(cfg.RootFolder())

	id := primitive.NewObjectID()

	file := &model.File{
		Id:       id,
		Filename: fmt.Sprintf("%s.txt", id.Hex()),
	}

	data := bytes.Buffer{}

	data.Write([]byte("cdsmklcdkmlkdsmckldsm"))

	err := ds.Save(file, data)

	assert.NoError(t, err)

	ds.Delete(file.Uri)
}
