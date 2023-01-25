package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"raspstore.github.io/file-manager/db"
	"raspstore.github.io/file-manager/model"
)

const filesCollectionName = "files"

type FilesRepository interface {
	Save(file *model.File) error
	FindById(id string) (*model.File, error)
	Delete(id string) error
	Update(file *model.File) error
	FindAll() (files []*model.File, err error)
}

type filesRepository struct {
	ctx  context.Context
	coll *mongo.Collection
}

func NewFilesRepository(ctx context.Context, conn db.MongoConnection) FilesRepository {
	return &filesRepository{ctx: ctx, coll: conn.DB().Collection(filesCollectionName)}
}

func (f *filesRepository) Save(file *model.File) error {
	if _, err := f.coll.InsertOne(f.ctx, file); err != nil {
		fmt.Println("could not create file metadata in database ", file, ", with error: ", err.Error())
		return err
	}

	return nil
}

func (f *filesRepository) FindById(id string) (user *model.File, err error) {
	var objectId primitive.ObjectID
	var file *model.File

	if value, err := primitive.ObjectIDFromHex(id); err != nil {
		return nil, err
	} else {
		objectId = value
	}

	if found := f.coll.FindOne(f.ctx, bson.M{"_id": objectId}); found.Err() != nil {
		return nil, found.Err()
	} else {
		err = found.Decode(&file)
		return file, err
	}
}

func (f *filesRepository) Delete(id string) error {
	var objectId primitive.ObjectID

	if value, err := primitive.ObjectIDFromHex(id); err != nil {
		return err
	} else {
		objectId = value
	}

	res, err := f.coll.DeleteOne(f.ctx, bson.M{"_id": objectId})

	if res.DeletedCount == 0 {
		return errors.New("file with provided id: " + id + " does not exists in database")
	}

	return err
}

func (f *filesRepository) Update(file *model.File) error {
	update := bson.M{"$set": bson.M{
		"filename":   file.Filename,
		"size":       file.Size,
		"updated_at": time.Now(),
		"updated_by": file.UpdatedBy}}

	f.coll.UpdateByID(f.ctx, file.Id, update)

	return nil
}

func (f *filesRepository) FindAll() (files []*model.File, err error) {
	var cursor *mongo.Cursor

	cursor, err = f.coll.Find(f.ctx, bson.M{})

	if err != nil {
		return nil, err
	}

	defer cursor.Close(f.ctx)

	for cursor.Next(f.ctx) {
		var file *model.File
		if err = cursor.Decode(&file); err != nil {
			return nil, err
		}

		files = append(files, file)
	}

	return files, nil
}
