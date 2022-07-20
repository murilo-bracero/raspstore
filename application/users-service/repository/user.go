package repository

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"raspstore.github.io/users-service/db"
	"raspstore.github.io/users-service/model"
	"raspstore.github.io/users-service/validators"
)

const usersCollectionName = "users"

type UsersRepository interface {
	Save(user *model.User) error
	FindById(id string) (user *model.User, err error)
	FindByEmail(email string) (user *model.User, err error)
	DeleteUser(id string) error
	UpdateUser(user *model.User) error
	FindAll() (users []*model.User, err error)
}

type usersRespository struct {
	ctx  context.Context
	coll *mongo.Collection
}

func NewUsersRepository(ctx context.Context, conn db.MongoConnection) UsersRepository {
	return &usersRespository{coll: conn.DB().Collection(usersCollectionName), ctx: ctx}
}

func (r *usersRespository) Save(user *model.User) error {
	user.UserId = uuid.NewString()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	_, err := r.coll.InsertOne(r.ctx, user)

	if err != nil {
		log.Println("Coud not create user ", user, " in MongoDB: ", err.Error())
		return err
	}

	return nil
}

func (r *usersRespository) FindById(id string) (user *model.User, err error) {

	res := r.coll.FindOne(r.ctx, bson.M{"user_id": id})

	if res.Err() == mongo.ErrNoDocuments {
		return nil, nil
	}

	err = res.Decode(&user)
	return user, err
}

func (r *usersRespository) FindByEmail(email string) (user *model.User, err error) {

	res := r.coll.FindOne(r.ctx, bson.M{"email": email})

	if res.Err() == mongo.ErrNoDocuments {
		return nil, nil
	}

	err = res.Decode(&user)
	return user, err
}

func (r *usersRespository) DeleteUser(id string) error {

	_, err := r.coll.DeleteOne(r.ctx, bson.M{"user_id": id})
	return err
}

func (r *usersRespository) UpdateUser(user *model.User) error {

	user.UpdatedAt = time.Now()

	res, err := r.coll.UpdateOne(r.ctx, bson.M{"user_id": user.UserId}, bson.M{
		"$set": bson.M{
			"username":     user.Username,
			"email":        user.Email,
			"phone_number": user.PhoneNumber,
			"updated_at":   user.UpdatedAt}})

	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return validators.ErrUserNotFound
	}

	if res.ModifiedCount == 0 {
		return errors.New("user could not be updated")
	}

	return nil
}

func (r *usersRespository) FindAll() (users []*model.User, err error) {
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
