package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
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
	FindByIdLookup(id string) (fileMetadata *model.FileMetadataLookup, err error)
	Delete(id string) error
	Update(file *model.File) error
	FindAll(page int, size int) (filesPage *model.FilePage, err error)
}

type filesRepository struct {
	ctx  context.Context
	coll *mongo.Collection
}

func NewFilesRepository(ctx context.Context, conn db.MongoConnection) FilesRepository {
	return &filesRepository{ctx: ctx, coll: conn.DB().Collection(filesCollectionName)}
}

func (f *filesRepository) Save(file *model.File) error {
	file.FileId = uuid.NewString()
	file.CreatedAt = time.Now()
	file.UpdatedAt = time.Now()

	if _, err := f.coll.InsertOne(f.ctx, file); err != nil {
		fmt.Println("could not create file metadata in database ", file, ", with error: ", err.Error())
		return err
	}

	return nil
}

func (f *filesRepository) FindById(id string) (file *model.File, err error) {
	found := f.coll.FindOne(f.ctx, bson.M{"file_id": id})

	err = found.Decode(&file)

	return file, err
}

func (f *filesRepository) FindByIdLookup(id string) (fileMetadata *model.FileMetadataLookup, err error) {
	match := bson.D{bson.E{Key: "$match", Value: bson.M{"file_id": id}}}

	pipeline := append([]bson.D{match}, lookupFileMetadata()...)

	cursor, err := f.coll.Aggregate(f.ctx, pipeline)

	if err != nil {
		return nil, err
	}

	for cursor.Next(f.ctx) {
		if err = cursor.Decode(&fileMetadata); err != nil {
			return nil, err
		}
	}

	return fileMetadata, nil
}

func (f *filesRepository) Delete(id string) error {
	_, err := f.coll.DeleteOne(f.ctx, bson.M{"file_id": id})

	return err
}

func (f *filesRepository) Update(file *model.File) error {
	filter := bson.M{"file_id": file.FileId}

	update := bson.M{"$set": bson.M{
		"filename":   file.Filename,
		"path":       file.Path,
		"editors":    file.Editors,
		"viewers":    file.Viewers,
		"updated_at": time.Now(),
		"updated_by": file.UpdatedBy}}

	f.coll.UpdateOne(f.ctx, filter, update)

	return nil
}

func (f *filesRepository) FindAll(page int, size int) (filesPage *model.FilePage, err error) {

	skip := bson.D{bson.E{Key: "$skip", Value: page * size}}

	limit := bson.D{bson.E{Key: "$limit", Value: size}}

	contentField := append([]bson.D{skip, limit}, lookupFileMetadata()...)

	totalCountField := []bson.M{{"$group": bson.M{"_id": nil, "count": bson.M{"$sum": 1}}}}
	facet := bson.D{
		primitive.E{Key: "$facet", Value: bson.D{
			primitive.E{Key: "content", Value: contentField}, primitive.E{Key: "totalCount", Value: totalCountField},
		}},
	}

	project := bson.D{
		primitive.E{Key: "$project", Value: bson.D{
			primitive.E{Key: "content", Value: "$content"},
			primitive.E{Key: "count", Value: bson.D{
				primitive.E{Key: "$arrayElemAt", Value: []interface{}{"$totalCount.count", 0}}}},
		}},
	}

	cursor, err := f.coll.Aggregate(f.ctx, mongo.Pipeline{facet, project})

	if err != nil {
		return nil, err
	}

	defer cursor.Close(f.ctx)

	for cursor.Next(f.ctx) {
		if err = cursor.Decode(&filesPage); err != nil {
			return nil, err
		}
	}

	return filesPage, nil
}

func lookupFileMetadata() []bson.D {
	lookupOwner := bson.D{bson.E{Key: "$lookup", Value: bson.M{"from": "users",
		"localField":   "owner_user_id",
		"foreignField": "user_id",
		"as":           "owner"}}}

	lookupCreatedBy := bson.D{bson.E{Key: "$lookup", Value: bson.M{"from": "users",
		"localField":   "created_by",
		"foreignField": "user_id",
		"as":           "created_by"}}}

	lookupUpdatedBy := bson.D{bson.E{Key: "$lookup", Value: bson.M{"from": "users",
		"localField":   "updated_by",
		"foreignField": "user_id",
		"as":           "updated_by"}}}

	lookupEditors := bson.D{bson.E{Key: "$lookup", Value: bson.M{"from": "users",
		"localField":   "editors",
		"foreignField": "user_id",
		"as":           "editors"}}}

	lookupViewers := bson.D{bson.E{Key: "$lookup", Value: bson.M{"from": "users",
		"localField":   "viewers",
		"foreignField": "user_id",
		"as":           "viewers"}}}

	unwindOwner := bson.D{bson.E{Key: "$unwind", Value: "$owner"}}
	unwindCreatedBy := bson.D{bson.E{Key: "$unwind", Value: "$created_by"}}
	unwindUpdatedBy := bson.D{bson.E{Key: "$unwind", Value: "$updated_by"}}

	return []bson.D{lookupOwner, unwindOwner, lookupCreatedBy, unwindCreatedBy, lookupUpdatedBy, unwindUpdatedBy, lookupEditors, lookupViewers}
}
