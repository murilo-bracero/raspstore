package repository

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"raspstore.github.io/authentication/db"
	"raspstore.github.io/authentication/model"
)

const credentialsCollectionName = "credentials"

type credentialsRespository struct {
	ctx  context.Context
	coll *mongo.Collection
}

func NewCredentialsRepository(ctx context.Context, conn db.MongoConnection) CredentialsRepository {
	return &credentialsRespository{coll: conn.DB().Collection(credentialsCollectionName), ctx: ctx}
}

func (r *credentialsRespository) Save(user *model.User, password string) error {

	hash, err := hash(password)

	if err != nil {
		return err
	}

	credential := model.Credential{
		Id:    user.UserId,
		Email: user.Email,
		Hash:  hash,
	}

	_, err = r.coll.InsertOne(r.ctx, credential)

	if err != nil {
		fmt.Println("Coud not create credentials ", user, " in MongoDB: ", err.Error())
		return err
	}

	return nil
}

func (r *credentialsRespository) Delete(id string) error {

	_, err := r.coll.DeleteOne(r.ctx, bson.M{"user_id": id})
	return err
}

func (r *credentialsRespository) Update(user *model.User) error {

	res, err := r.coll.UpdateOne(r.ctx, bson.M{"user_id": user.UserId}, bson.M{
		"$set": bson.M{
			"email": user.Email}})

	if err == nil && (res.MatchedCount == 0 || res.ModifiedCount == 0) {
		return errors.New("credential could not be updated")
	}

	return err
}

func (r *credentialsRespository) IsCredentialsCorrect(email string, password string) bool {
	var credential model.Credential

	res := r.coll.FindOne(r.ctx, bson.M{"email": email})

	if res.Err() != nil {
		return false
	}

	if err := res.Decode(&credential); err != nil {
		return false
	}

	return bcrypt.CompareHashAndPassword([]byte(credential.Hash), []byte(password)) == nil
}

func hash(text string) (hash string, err error) {
	bts := []byte(text)

	raw, err := bcrypt.GenerateFromPassword(bts, bcrypt.DefaultCost)
	return string(raw[:]), err
}
