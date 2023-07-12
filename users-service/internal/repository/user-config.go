package repository

import (
	"context"

	"github.com/murilo-bracero/raspstore/users-service/internal/database"
	"github.com/murilo-bracero/raspstore/users-service/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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

func (r *usersConfigRepository) Update(usersConfig *model.UserConfiguration) error {

	_, err := r.coll.UpdateOne(r.ctx, bson.M{}, bson.M{
		"$set": bson.M{
			"storage_limit":              usersConfig.StorageLimit,
			"min_password_length":        usersConfig.MinPasswordLength,
			"allow_public_user_creation": usersConfig.AllowPublicUserCreation,
			"allow_login_with_email":     usersConfig.AllowLoginWithEmail,
			"enforce_mfa":                usersConfig.EnforceMfa}})

	return err
}
