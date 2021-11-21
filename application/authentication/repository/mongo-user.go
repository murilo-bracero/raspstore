package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"raspstore.github.io/authentication/db"
	"raspstore.github.io/authentication/model"
)

const UsersCollectionName = "users"

type mongoUsersRespository struct {
	ctx  context.Context
	coll *mongo.Collection
}

func NewMongoUsersRepository(ctx context.Context, conn db.MongoConnection) UsersRepository {
	return &mongoUsersRespository{coll: conn.DB().Collection(UsersCollectionName), ctx: ctx}
}

func (r *mongoUsersRespository) Save(user *model.User) error {
	user.Id = primitive.NewObjectID()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	_, err := r.coll.InsertOne(r.ctx, user)

	if err != nil {
		fmt.Println("Coud not create user ", user, " in MongoDB: ", err.Error())
		return err
	}

	return nil
}

func (r *mongoUsersRespository) FindById(id string) (user *model.User, err error) {
	objectId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		fmt.Println("Could not convert provided ID: ", id, " to a valid ObjectId: ", err.Error())
		return nil, err
	}

	res := r.coll.FindOne(r.ctx, bson.M{"_id": objectId})
	err = res.Decode(&user)
	return user, err
}

func (r *mongoUsersRespository) FindByEmailOrUsername(email string, username string) (user *model.User, err error) {

	res := r.coll.FindOne(r.ctx, bson.M{"$or": [2]bson.M{{"email": email}, {"username": username}}})
	err = res.Decode(&user)
	return user, err
}

func (r *mongoUsersRespository) DeleteUser(id string) error {

	objectId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return err
	}

	_, err = r.coll.DeleteOne(r.ctx, bson.M{"_id": objectId})
	return err
}

func (r *mongoUsersRespository) UpdateUser(user *model.User) error {

	user.UpdatedAt = time.Now()

	res, err := r.coll.UpdateOne(r.ctx, bson.M{"_id": user.Id}, bson.M{
		"$set": bson.M{
			"username":     user.Username,
			"email":        user.Email,
			"phone_number": user.PhoneNumber,
			"updated_at":   user.UpdatedAt}})

	if err == nil && (res.MatchedCount == 0 || res.ModifiedCount == 0) {
		return errors.New("user could not be updated")
	}

	return err
}

func (r *mongoUsersRespository) FindAll() (users []*model.User, err error) {
	var cursor *mongo.Cursor

	cursor, err = r.coll.Find(r.ctx, bson.M{})

	if err != nil {
		return nil, err
	}

	defer cursor.Close(r.ctx)

	for cursor.Next(r.ctx) {
		var user *model.User
		if err = cursor.Decode(&user); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}
