package repository

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"raspstore.github.io/users-service/internal"
	"raspstore.github.io/users-service/internal/database"
	"raspstore.github.io/users-service/internal/model"
)

const usersCollectionName = "users"

type UsersRepository interface {
	Save(user *model.User) error
	ExistsByEmailOrUsername(email string, username string) (bool, error)
	FindById(id string) (user *model.User, err error)
	FindByEmail(email string) (user *model.User, err error)
	FindAll(page int, size int) (userPage *model.UserPage, err error)
	Update(user *model.User) error
	Delete(id string) error
}

type usersRespository struct {
	ctx  context.Context
	coll *mongo.Collection
}

func NewUsersRepository(ctx context.Context, conn database.MongoConnection) UsersRepository {
	return &usersRespository{coll: conn.Collection(usersCollectionName), ctx: ctx}
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

func (r *usersRespository) ExistsByEmailOrUsername(email string, username string) (bool, error) {
	filter := bson.M{"$or": []bson.M{{"email": email}, {"username": username}}}

	count, err := r.coll.CountDocuments(r.ctx, filter)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *usersRespository) FindById(id string) (user *model.User, err error) {
	res := r.coll.FindOne(r.ctx, bson.M{"user_id": id})

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

func (r *usersRespository) FindAll(page int, size int) (userPage *model.UserPage, err error) {
	skip := page * size

	contentField := []bson.M{{"$skip": skip}, {"$limit": size}}
	totalCountField := []bson.M{{"$group": bson.M{"_id": nil, "count": bson.M{"$sum": 1}}}}
	facet := bson.D{
		primitive.E{Key: "$facet", Value: bson.D{
			primitive.E{Key: "content", Value: contentField}, primitive.E{Key: "totalCount", Value: totalCountField},
		}},
	}

	project := bson.D{
		primitive.E{Key: "$project", Value: bson.D{
			primitive.E{Key: "content", Value: "$content"},
			primitive.E{Key: "count", Value: bson.D{
				primitive.E{Key: "$arrayElemAt", Value: []interface{}{"$totalCount.count", 0}}}},
		}},
	}

	cursor, err := r.coll.Aggregate(r.ctx, mongo.Pipeline{facet, project}, &options.AggregateOptions{})

	if err != nil {
		return nil, err
	}

	defer cursor.Close(r.ctx)

	for cursor.Next(r.ctx) {
		if err = cursor.Decode(&userPage); err != nil {
			return nil, err
		}
	}

	return userPage, nil
}

func (r *usersRespository) Update(user *model.User) error {
	user.UpdatedAt = time.Now()

	res, err := r.coll.UpdateOne(r.ctx, bson.M{"user_id": user.UserId}, bson.M{
		"$set": bson.M{
			"username":     user.Username,
			"email":        user.Email,
			"password":     user.PasswordHash,
			"phone_number": user.PhoneNumber,
			"updated_at":   user.UpdatedAt}})

	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return internal.ErrUserNotFound
	}

	if res.ModifiedCount == 0 {
		return errors.New("user could not be updated")
	}

	return nil
}

func (r *usersRespository) Delete(id string) error {
	_, err := r.coll.DeleteOne(r.ctx, bson.M{"user_id": id})
	return err
}
