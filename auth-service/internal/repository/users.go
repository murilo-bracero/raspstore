package repository

import (
	"context"
	"errors"

	"github.com/murilo-bracero/raspstore/auth-service/internal/database"
	"github.com/murilo-bracero/raspstore/auth-service/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const usersCollectionName = "users"

type UsersRepository interface {
	Update(usr *model.User) error
	FindByUsername(username string) (usr *model.User, err error)
}

type usersRespository struct {
	ctx  context.Context
	coll *mongo.Collection
}

func NewUsersRepository(ctx context.Context, conn database.MongoConnection) UsersRepository {
	return &usersRespository{coll: conn.Collection(usersCollectionName), ctx: ctx}
}

func (r *usersRespository) FindByUsername(username string) (usr *model.User, err error) {
	res := r.coll.FindOne(r.ctx, bson.M{"username": username})

	err = res.Decode(&usr)
	return usr, err
}

func (r *usersRespository) Update(usr *model.User) error {

	filter := bson.D{{Key: "user_id", Value: usr.UserId}}

	update := bson.D{{Key: "$set",
		Value: bson.D{
			{Key: "is_mfa_enabled", Value: usr.IsMfaEnabled},
			{Key: "is_mfa_verified", Value: usr.IsMfaVerified},
			{Key: "secret", Value: usr.Secret},
			{Key: "refresh_token", Value: usr.RefreshToken},
		},
	}}

	res, err := r.coll.UpdateOne(r.ctx, filter, update)

	if err == nil && (res.MatchedCount == 0 || res.ModifiedCount == 0) {
		return errors.New("credential could not be updated")
	}

	return err
}
