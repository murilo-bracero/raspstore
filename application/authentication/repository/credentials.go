package repository

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"raspstore.github.io/authentication/db"
	"raspstore.github.io/authentication/model"
)

const credentialsCollectionName = "credentials"

type mongoCredentialsRespository struct {
	ctx  context.Context
	coll *mongo.Collection
}

func NewMongoCredentialsRepository(ctx context.Context, conn db.MongoConnection) CredentialsRepository {
	return &mongoCredentialsRespository{coll: conn.DB().Collection(credentialsCollectionName), ctx: ctx}
}

func (r *mongoCredentialsRespository) Save(user *model.User, password string) error {

	hash, err := hash(password)

	if err != nil {
		return err
	}

	document := bson.M{
		"_id":      primitive.NewObjectID(),
		"user_id":  user.UserId,
		"email":    user.Email,
		"password": hash,
	}

	_, err = r.coll.InsertOne(r.ctx, document)

	if err != nil {
		fmt.Println("Coud not create credentials ", user, " in MongoDB: ", err.Error())
		return err
	}

	return nil
}

func (r *mongoCredentialsRespository) Delete(id string) error {

	_, err := r.coll.DeleteOne(r.ctx, bson.M{"user_id": id})
	return err
}

func (r *mongoCredentialsRespository) Update(user *model.User) error {

	res, err := r.coll.UpdateOne(r.ctx, bson.M{"user_id": user.UserId}, bson.M{
		"$set": bson.M{
			"email": user.Email}})

	if err == nil && (res.MatchedCount == 0 || res.ModifiedCount == 0) {
		return errors.New("credential could not be updated")
	}

	return err
}

func (r *mongoCredentialsRespository) Authenticate(token string) (uid string, err error) {

	secret := os.Getenv("JWT_SECRET")

	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("error reading jwt: wrong signing method: %v", t.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return "", err
	}

	return parsedToken.Claims.(jwt.MapClaims)["uid"].(string), nil
}

func hash(text string) (hash string, err error) {
	bts := []byte(text)

	raw, err := bcrypt.GenerateFromPassword(bts, bcrypt.DefaultCost)
	return string(raw[:]), err
}
