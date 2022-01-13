package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/pquerna/otp/totp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"raspstore.github.io/authentication/db"
	"raspstore.github.io/authentication/model"
)

const credentialsCollectionName = "credentials"

type CredentialsRepository interface {
	Save(cred *model.Credential) error
	Update(cred *model.Credential) error
	Delete(id string) error
	IsCredentialsCorrect(email string, password string, token string) bool
	Has2FAEnabledByEmail(email string) bool
	Has2FAEnabledByUID(uid string) bool
	FindById(id string) (cred *model.Credential, err error)
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

func (r *credentialsRespository) FindById(id string) (cred *model.Credential, err error) {
	res := r.coll.FindOne(r.ctx, bson.M{"user_id": id})

	err = res.Decode(&cred)
	return cred, err
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

func (r *credentialsRespository) Has2FAEnabledByEmail(email string) bool {
	var credential model.Credential

	res := r.coll.FindOne(r.ctx, bson.M{"email": email})

	if res.Err() != nil {
		return false
	}

	if err := res.Decode(&credential); err != nil {
		return false
	}

	return credential.Has2FAEnabled
}

func (r *credentialsRespository) Has2FAEnabledByUID(uid string) bool {
	var credential model.Credential

	res := r.coll.FindOne(r.ctx, bson.M{"user_id": uid})

	if res.Err() != nil {
		return false
	}

	if err := res.Decode(&credential); err != nil {
		return false
	}

	return credential.Has2FAEnabled
}

func (r *credentialsRespository) IsCredentialsCorrect(email string, password string, token string) bool {
	var credential model.Credential

	res := r.coll.FindOne(r.ctx, bson.M{"email": email})

	if res.Err() != nil {
		return false
	}

	if err := res.Decode(&credential); err != nil {
		return false
	}

	if bcrypt.CompareHashAndPassword([]byte(credential.Hash), []byte(password)) == nil {
		if credential.Has2FAEnabled {
			return totp.Validate(token, credential.Secret)
		}

		return true
	}

	return false
}
