package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"raspstore.github.io/users-service/internal/database"
	"raspstore.github.io/users-service/internal/model"
)

const usersConfigCollection = "users-config"

type UsersConfigRepository interface {
	Find() (usersConfig *model.UserConfiguration, err error)
	Update(usersConfig *model.UserConfiguration) error
}

type usersConfigRepository struct {
	ctx  context.Context
	coll *mongo.Collection
}

func NewUsersConfigRepository(ctx context.Context, conn database.MongoConnection) UsersConfigRepository {
	return &usersConfigRepository{coll: conn.Collection(usersConfigCollection), ctx: ctx}
}

func (r *usersConfigRepository) Find() (usersConfig *model.UserConfiguration, err error) {
	result := r.coll.FindOne(r.ctx, bson.M{})

	err = result.Decode(&usersConfig)

	return usersConfig, err
}

//TODO
func (r *usersConfigRepository) Update(usersConfig *model.UserConfiguration) error {
	return nil
}
