package repository

import (
	"context"
	"time"

	"github.com/murilo-bracero/raspstore/file-service/internal"
	db "github.com/murilo-bracero/raspstore/file-service/internal/database"
	"github.com/murilo-bracero/raspstore/file-service/internal/model"
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
	FindAll(userId string, page int, size int, filename string, secret bool) (filesPage *model.FilePage, err error)
}

type filesRepository struct {
	ctx  context.Context
	coll *mongo.Collection
}

func NewFilesRepository(ctx context.Context, conn db.MongoConnection) FilesRepository {
	return &filesRepository{ctx: ctx, coll: conn.Collection(filesCollectionName)}
}

func (f *filesRepository) Save(file *model.File) error {
	file.CreatedAt = time.Now()
	file.UpdatedAt = time.Now()

	if _, err := f.coll.InsertOne(f.ctx, file); err != nil {
		return err
	}

	return nil
}

func (f *filesRepository) FindById(userId string, fileId string) (file *model.File, err error) {
	filter := bson.D{
		bson.E{Key: "file_id", Value: fileId},
		bson.E{Key: "$or",
			Value: bson.A{
				bson.D{bson.E{Key: "owner_user_id", Value: userId}},
				bson.D{bson.E{Key: "editors", Value: userId}},
				bson.D{bson.E{Key: "viewers", Value: userId}},
			},
		},
	}

	found := f.coll.FindOne(f.ctx, filter)

	if found.Err() == mongo.ErrNoDocuments {
		return nil, internal.ErrFileDoesNotExists
	}

	err = found.Decode(&file)

	return file, err
}

func (f *filesRepository) FindByIdLookup(userId string, fileId string) (fileMetadata *model.FileMetadataLookup, err error) {
	match := bson.D{
		bson.E{Key: "$match", Value: bson.D{
			bson.E{Key: "file_id", Value: fileId},
		}},
	}

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
	filter := bson.D{
		bson.E{Key: "file_id", Value: fileId},
		bson.E{Key: "owner_user_id", Value: userId},
	}

	_, err := f.coll.DeleteOne(f.ctx, filter)

	return err
}

func (f *filesRepository) Update(userId string, file *model.File) error {

	var filter bson.D

	if file.Secret {
		filter = bson.D{
			bson.E{Key: "file_id", Value: file.FileId},
			bson.E{Key: "$or",
				Value: bson.A{
					bson.D{bson.E{Key: "owner_user_id", Value: userId}},
					bson.D{bson.E{Key: "editors", Value: userId}},
				},
			},
		}
	} else {
		filter = bson.D{
			bson.E{Key: "file_id", Value: file.FileId},
			bson.E{Key: "owner_user_id", Value: userId},
		}
	}

	update := bson.M{"$set": bson.M{
		"filename":   file.Filename,
		"is_secret":  file.Secret,
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

func (f *filesRepository) FindAll(userId string, page int, size int, filename string, secret bool) (filesPage *model.FilePage, err error) {
	contentField := []bson.D{}

	if filename != "" {
		contentField = append(contentField, bson.D{bson.E{Key: "$match", Value: bson.D{
			bson.E{Key: "filename", Value: bson.D{
				bson.E{Key: "$regex", Value: filename},
			}},
		}}})
	}

	var accessControl bson.D
	if secret {
		accessControl = bson.D{bson.E{Key: "$match", Value: bson.D{bson.E{Key: "owner_user_id", Value: userId}}}}
		contentField = append(contentField, bson.D{bson.E{Key: "$match", Value: bson.D{
			bson.E{Key: "is_secret", Value: true},
		}}})
	} else {
		accessControl = aggregateAccessControl(userId)
	}

	contentField = append(contentField, bson.D{bson.E{Key: "$skip", Value: page * size}}, bson.D{bson.E{Key: "$limit", Value: size}})

	contentField = append(contentField, lookupUserFields()...)

	totalCountField := bson.D{
		bson.E{Key: "$group", Value: bson.D{
			bson.E{Key: "_id", Value: nil},
			bson.E{Key: "count", Value: bson.D{
				bson.E{Key: "$sum", Value: 1},
			}},
		}},
	}

	facet := bson.D{
		bson.E{Key: "$facet", Value: bson.D{
			bson.E{Key: "content", Value: contentField},
			bson.E{Key: "totalCount", Value: bson.A{totalCountField}},
		}},
	}

	project := bson.D{
		bson.E{Key: "$project", Value: bson.D{
			bson.E{Key: "content", Value: "$content"},
			bson.E{Key: "count", Value: bson.D{
				bson.E{Key: "$arrayElemAt", Value: bson.A{"$totalCount.count", 0}},
			}}}},
	}

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
	match := bson.D{
		bson.E{Key: "$match", Value: bson.D{
			bson.E{Key: "owner_user_id", Value: userId},
		}},
	}

	project := bson.D{
		primitive.E{Key: "$group", Value: bson.M{
			"_id":        "$owner_user_id",
			"totalUsage": bson.M{"$sum": "$size"},
		}},
	}

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
			return -1, nil
		}
	}

	return
}

func aggregateAccessControl(userId string) bson.D {
	return bson.D{bson.E{Key: "$match", Value: bson.D{
		bson.E{Key: "$or", Value: bson.A{
			bson.D{bson.E{Key: "owner_user_id", Value: userId}},
			bson.D{bson.E{Key: "viewers", Value: userId}},
			bson.D{bson.E{Key: "editors", Value: userId}},
		}}}}}
}

func lookupUserFields() []bson.D {
	return []bson.D{
		lookup("users", "owner_user_id", "user_id", "owner"),
		lookup("users", "created_by", "user_id", "created_by"),
		lookup("users", "updated_by", "user_id", "updated_by"),
		lookup("users", "viewers", "user_id", "viewers"),
		lookup("users", "editors", "user_id", "editors"),
		{bson.E{Key: "$unwind", Value: "$owner"}},
		{bson.E{Key: "$unwind", Value: "$created_by"}},
		{bson.E{Key: "$unwind", Value: "$updated_by"}},
	}
}

func lookup(from string, lField string, fField string, as string) bson.D {
	return bson.D{bson.E{Key: "$lookup", Value: bson.D{
		bson.E{Key: "from", Value: from},
		bson.E{Key: "localField", Value: lField},
		bson.E{Key: "foreignField", Value: fField},
		bson.E{Key: "as", Value: as},
	}}}
}
