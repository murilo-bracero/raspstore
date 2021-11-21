package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"raspstore.github.io/authentication/pb"
)

type User struct {
	Id          primitive.ObjectID `bson:"_id" datastore:"-"`
	UserId      string             `bson:"-" datastore:"user_id"`
	Username    string             `bson:"username" datastore:"username"`
	Email       string             `bson:"email" datastore:"email"`
	PhoneNumber string             `bson:"phone_number" datastore:"phone_number"`
	CreatedAt   time.Time          `bson:"created_at" datastore:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" datastore:"updated_at"`
}

func (u *User) ToProtoBuffer() *pb.User {
	proto := &pb.User{
		Username:    u.Username,
		Email:       u.Email,
		PhoneNumber: u.PhoneNumber,
		CreatedAt:   u.CreatedAt.Unix(),
		UpdatedAt:   u.UpdatedAt.Unix(),
	}

	if u.Id == primitive.NilObjectID {
		proto.Id = u.UserId
	} else {
		proto.Id = u.Id.Hex()
	}

	return proto
}

func (u *User) FromProtoBuffer(user *pb.CreateUserRequest) {

	u.Username = user.Username
	u.Email = user.Email
	u.PhoneNumber = user.PhoneNumber
}

func (u *User) FromUpdateProto(user *pb.UpdateUserRequest) error {

	if len(user.Id) > 24 {
		u.UserId = user.Id
	} else {
		oid, err := primitive.ObjectIDFromHex(user.Id)

		if err != nil {
			return err
		}

		u.Id = oid
	}

	u.Username = user.Username
	u.Email = user.Email
	u.PhoneNumber = user.PhoneNumber
	return nil
}
