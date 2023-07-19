package repository

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/murilo-bracero/raspstore/file-info-service/internal"
	db "github.com/murilo-bracero/raspstore/file-info-service/internal/database"
	"github.com/murilo-bracero/raspstore/file-info-service/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const filesCollectionName = "files"

type FilesRepository interface {
	Save(file *model.File) error
	FindById(userId string, fileId string) (*model.File, error)
	FindByIdLookup(userId string, fileId string) (fileMetadata *model.FileMetadataLookup, err error)
	FindUsageByUserId(userId string) (usage int64, err error)
	Delete(userId string, fileId string) error
	Update(userId string, file *model.File) error
	FindAll(userId string, page int, size int) (filesPage *model.FilePage, err error)
}

type filesRepository struct {
	ctx  context.Context
	coll *mongo.Collection
}

func NewFilesRepository(ctx context.Context, conn db.MongoConnection) FilesRepository {
	return &filesRepository{ctx: ctx, coll: conn.Collection(filesCollectionName)}
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

func (f *filesRepository) FindById(userId string, fileId string) (file *model.File, err error) {
	found := f.coll.FindOne(f.ctx, bson.M{"file_id": fileId, "$or": addAnyPermissionFilter(userId)})

	if found.Err() == mongo.ErrNoDocuments {
		return nil, internal.ErrFileDoesNotExists
	}

	err = found.Decode(&file)

	return file, err
}

func (f *filesRepository) FindByIdLookup(userId string, fileId string) (fileMetadata *model.FileMetadataLookup, err error) {
	match := bson.D{bson.E{Key: "$match", Value: filterByFileId(fileId)}}

	pipeline := append([]bson.D{match}, lookupUserFields()...)

	pipeline = append(pipeline, aggregateAccessControl(userId))

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

func (f *filesRepository) Delete(userId string, fileId string) error {
	filter := filterByFileId(fileId)
	addOwnerPermissionFilter(filter, userId)

	_, err := f.coll.DeleteOne(f.ctx, filter)

	return err
}

func (f *filesRepository) Update(userId string, file *model.File) error {
	filter := filterByFileId(file.FileId)
	addEditorPermissionFilter(filter, userId)

	update := bson.M{"$set": bson.M{
		"filename":   file.Filename,
		"path":       file.Path,
		"editors":    file.Editors,
		"viewers":    file.Viewers,
		"updated_at": time.Now(),
		"updated_by": file.UpdatedBy}}

	result, err := f.coll.UpdateOne(f.ctx, filter, update)

	if result.MatchedCount == 0 {
		return internal.ErrFileDoesNotExists
	}

	return err
}

func (f *filesRepository) FindAll(userId string, page int, size int) (filesPage *model.FilePage, err error) {

	skip := bson.D{bson.E{Key: "$skip", Value: page * size}}

	limit := bson.D{bson.E{Key: "$limit", Value: size}}

	contentField := append([]bson.D{skip, limit}, lookupUserFields()...)

	totalCountField := []bson.M{{"$group": bson.M{"_id": nil, "count": bson.M{"$sum": 1}}}}

	accessControl := aggregateAccessControl(userId)
	facet := facet(contentField, totalCountField)
	project := projectPage()

	cursor, err := f.coll.Aggregate(f.ctx, mongo.Pipeline{accessControl, facet, project})

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

func (f *filesRepository) FindUsageByUserId(userId string) (usage int64, err error) {
	match := bson.D{bson.E{Key: "$match", Value: filterByOwnerId(userId)}}

	project := groupUserUsage()

	cursor, err := f.coll.Aggregate(f.ctx, mongo.Pipeline{match, project})

	if err != nil {
		return 0, err
	}

	defer cursor.Close(f.ctx)

	for cursor.Next(f.ctx) {
		value, err := cursor.Current.LookupErr("totalUsage")

		if err != nil {
			return 0, nil
		}

		var ok bool
		if usage, ok = value.Int64OK(); !ok {
			log.Println("[WARN] Could not convert usage into a valid duble value")
			return 0, nil
		}
	}

	return
}

func filterByFileId(fileId string) bson.M {
	return bson.M{"file_id": fileId}
}

func filterByOwnerId(userId string) bson.M {
	return bson.M{"owner_user_id": userId}
}

func addAnyPermissionFilter(userId string) []bson.M {
	ownerClause := filterByOwnerId(userId)

	viewerClause := bson.M{"viewers": userId}

	editorClause := bson.M{"editors": userId}

	return []bson.M{ownerClause, viewerClause, editorClause}
}

func addEditorPermissionFilter(query bson.M, userId string) {
	query["editors"] = userId
}

func addOwnerPermissionFilter(query bson.M, userId string) {
	query["owner_user_id"] = userId
}

func facet(contentField []primitive.D, totalCountField []primitive.M) bson.D {
	return bson.D{
		primitive.E{Key: "$facet", Value: bson.D{
			primitive.E{Key: "content", Value: contentField}, primitive.E{Key: "totalCount", Value: totalCountField},
		}},
	}
}

func groupUserUsage() bson.D {
	return bson.D{
		primitive.E{Key: "$group", Value: bson.M{
			"_id":        "$owner_user_id",
			"totalUsage": bson.M{"$sum": "$size"},
		}},
	}
}

func projectPage() bson.D {
	return bson.D{
		primitive.E{Key: "$project", Value: bson.D{
			primitive.E{Key: "content", Value: "$content"},
			primitive.E{Key: "count", Value: bson.D{
				primitive.E{Key: "$arrayElemAt", Value: []interface{}{"$totalCount.count", 0}}}},
		}},
	}
}

func aggregateAccessControl(userId string) bson.D {
	orClauses := addAnyPermissionFilter(userId)

	orClause := bson.D{bson.E{Key: "$or", Value: orClauses}}

	return bson.D{bson.E{Key: "$match", Value: orClause}}
}

func lookupUserFields() []bson.D {
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
