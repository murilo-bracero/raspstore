package repository

import (
	"context"
	"time"

	"github.com/murilo-bracero/raspstore/auth-service/internal"
	"github.com/murilo-bracero/raspstore/auth-service/internal/database"
	"github.com/murilo-bracero/raspstore/auth-service/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const usersCollectionName = "users"

type UsersRepository interface {
	Update(usr *model.User) error
	FindByUsername(username string) (usr *model.User, err error)
	FindByUserId(userId string) (user *model.User, err error)
	ExistsByUsername(username string) (bool, error)
	Delete(userId string) error
}

type usersRespository struct {
	ctx  context.Context
	coll *mongo.Collection
}

func NewUsersRepository(ctx context.Context, conn database.MongoConnection) UsersRepository {
	return &usersRespository{coll: conn.Collection(usersCollectionName), ctx: ctx}
}

func (r *usersRespository) FindByUserId(userId string) (user *model.User, err error) {
	res := r.coll.FindOne(r.ctx, bson.D{{Key: "user_id", Value: userId}})

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

func (r *usersRespository) Update(usr *model.User) error {
	usr.UpdatedAt = time.Now()

	filter := bson.D{{Key: "user_id", Value: usr.UserId}}

	update := bson.D{{Key: "$set",
		Value: bson.D{
			{Key: "username", Value: usr.Username},
			{Key: "is_mfa_enabled", Value: usr.IsMfaEnabled},
			{Key: "is_mfa_verified", Value: usr.IsMfaVerified},
			{Key: "secret", Value: usr.Secret},
			{Key: "refresh_token", Value: usr.RefreshToken},
			{Key: "updated_at", Value: usr.UpdatedAt},
		},
	}}

	_, err := r.coll.UpdateOne(r.ctx, filter, update)

	return err
}

func (r *usersRespository) ExistsByUsername(username string) (bool, error) {
	filter := bson.M{"username": username}

	count, err := r.coll.CountDocuments(r.ctx, filter)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *usersRespository) Delete(userId string) error {
	filter := bson.D{{Key: "user_id", Value: userId}}

	_, err := r.coll.DeleteOne(r.ctx, filter)

	return err
}
