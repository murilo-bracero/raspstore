package repository

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"raspstore.github.io/users-service/db"
	"raspstore.github.io/users-service/model"
)

const credentialsCollectionName = "credentials"

type CredentialsRepository interface {
	Save(cred *model.Credential) error
	Update(cred *model.Credential) error
	Delete(id string) error
}

type credentialsRespository struct {
	ctx  context.Context
	coll *mongo.Collection
}

func NewCredentialsRepository(ctx context.Context, conn db.MongoConnection) CredentialsRepository {
	return &credentialsRespository{coll: conn.DB().Collection(credentialsCollectionName), ctx: ctx}
}

func (r *credentialsRespository) Save(cred *model.Credential) error {

	if _, err := r.coll.InsertOne(r.ctx, cred); err != nil {
		fmt.Println("Coud not create credentials ", cred, " in MongoDB: ", err.Error())
		return err
	}

	return nil
}

func (r *credentialsRespository) Delete(id string) error {

	_, err := r.coll.DeleteOne(r.ctx, bson.M{"user_id": id})
	return err
}

func (r *credentialsRespository) Update(cred *model.Credential) error {

	filter := bson.D{{Key: "user_id", Value: cred.Id}}

	update := bson.D{{Key: "$set",
		Value: bson.D{
			{Key: "email", Value: cred.Email},
			{Key: "has_2FA_enabled", Value: cred.Has2FAEnabled},
			{Key: "password", Value: cred.Hash},
			{Key: "secret", Value: cred.Secret},
		},
	}}

	res, err := r.coll.UpdateOne(r.ctx, filter, update)

	if err == nil && (res.MatchedCount == 0 || res.ModifiedCount == 0) {
		return errors.New("credential could not be updated")
	}

	return err
}
