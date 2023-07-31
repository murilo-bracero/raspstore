package repository

import (
	"context"
	"time"

	"github.com/murilo-bracero/raspstore/idp/internal"
	"github.com/murilo-bracero/raspstore/idp/internal/database"
	"github.com/murilo-bracero/raspstore/idp/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const usersCollectionName = "users"

type UsersRepository interface {
	Save(usr *model.User) error
	Update(usr *model.User) error
	FindAll(page int, size int, username string, enabled *bool) (userPage *model.UserPage, err error)
	FindByUsername(username string) (usr *model.User, err error)
	FindByUserId(userId string) (user *model.User, err error)
	Delete(userId string) error
}

type usersRespository struct {
	ctx  context.Context
	coll *mongo.Collection
}

func NewUsersRepository(ctx context.Context, conn database.MongoConnection) UsersRepository {
	return &usersRespository{coll: conn.Collection(usersCollectionName), ctx: ctx}
}

func (r *usersRespository) Save(usr *model.User) error {
	usr.BeforeCreate()

	_, err := r.coll.InsertOne(r.ctx, usr)

	if mongo.IsDuplicateKeyError(err) {
		return internal.ErrConflict
	}

	if err != nil {
		return err
	}

	return nil
}

func (r *usersRespository) Update(usr *model.User) error {
	usr.UpdatedAt = time.Now()

	filter := bson.D{{Key: "user_id", Value: usr.UserId}}

	update := bson.D{{Key: "$set",
		Value: bson.D{
			{Key: "username", Value: usr.Username},
			{Key: "is_enabled", Value: usr.IsEnabled},
			{Key: "permissions", Value: usr.Permissions},
			{Key: "is_mfa_enabled", Value: usr.IsMfaEnabled},
			{Key: "is_mfa_verified", Value: usr.IsMfaVerified},
			{Key: "secret", Value: usr.Secret},
			{Key: "refresh_token", Value: usr.RefreshToken},
			{Key: "updated_at", Value: usr.UpdatedAt},
		},
	}}

	_, err := r.coll.UpdateOne(r.ctx, filter, update)

	if mongo.IsDuplicateKeyError(err) {
		return internal.ErrConflict
	}

	return err
}

func (r *usersRespository) FindAll(page int, size int, username string, enabled *bool) (userPage *model.UserPage, err error) {
	contentField := []bson.D{}

	if username != "" {
		contentField = append(contentField, bson.D{{Key: "$match", Value: bson.D{
			{Key: "username", Value: bson.D{
				{Key: "$regex", Value: username},
			}},
		}}})
	}

	if enabled != nil {
		contentField = append(contentField, bson.D{{Key: "$match", Value: bson.D{
			{Key: "is_enabled", Value: *enabled},
		}}})
	}

	contentField = append(contentField, bson.D{{Key: "$skip", Value: page * size}}, bson.D{{Key: "$limit", Value: size}})

	totalCountField := []bson.M{{"$group": bson.M{"_id": nil, "count": bson.M{"$sum": 1}}}}
	facet := bson.D{
		{Key: "$facet", Value: bson.D{
			{Key: "content", Value: contentField}, {Key: "totalCount", Value: totalCountField},
		}},
	}

	project := bson.D{
		{Key: "$project", Value: bson.D{
			{Key: "content", Value: "$content"},
			{Key: "count", Value: bson.D{
				{Key: "$arrayElemAt", Value: bson.A{"$totalCount.count", 0}}}},
		}},
	}

	cursor, err := r.coll.Aggregate(r.ctx, mongo.Pipeline{facet, project})

	if err != nil {
		return nil, err
	}

	defer cursor.Close(r.ctx)

	for cursor.Next(r.ctx) {
		if err = cursor.Decode(&userPage); err != nil {
			return nil, err
		}
	}

	return userPage, nil
}

func (r *usersRespository) FindByUserId(userId string) (user *model.User, err error) {
	res := r.coll.FindOne(r.ctx, bson.D{{Key: "user_id", Value: userId}})

	if res.Err() == mongo.ErrNoDocuments {
		return nil, internal.ErrUserNotFound
	}

	err = res.Decode(&user)
	return
}

func (r *usersRespository) FindByUsername(username string) (usr *model.User, err error) {
	res := r.coll.FindOne(r.ctx, bson.M{"username": username})

	err = res.Decode(&usr)

	if err == mongo.ErrNoDocuments {
		return nil, internal.ErrUserNotFound
	}

	return usr, err
}

func (r *usersRespository) Delete(userId string) error {
	filter := bson.D{{Key: "user_id", Value: userId}}

	_, err := r.coll.DeleteOne(r.ctx, filter)

	return err
}
