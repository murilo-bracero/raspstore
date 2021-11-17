package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"raspstore.github.io/authentication/pb"
)

type User struct {
	Id          primitive.ObjectID `bson:"_id"`
	Username    string             `bson:"username"`
	Email       string             `bson:"email"`
	PhoneNumber string             `bson:"phone_number"`
	CreatedAt   time.Time          `bson:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at"`
}

func (u *User) ToProtoBuffer() *pb.User {
	return &pb.User{
		Id:          u.Id.Hex(),
		Username:    u.Username,
		Email:       u.Email,
		PhoneNumber: u.PhoneNumber,
		CreatedAt:   u.CreatedAt.Unix(),
		UpdatedAt:   u.UpdatedAt.Unix(),
	}
}

func (u *User) FromProtoBuffer(user *pb.CreateUserRequest) {

	u.Username = user.Username
	u.Email = user.Email
	u.PhoneNumber = user.PhoneNumber
}

func (u *User) FromUpdateProto(user *pb.UpdateUserRequest) error {
	oid, err := primitive.ObjectIDFromHex(user.Id)

	if err != nil {
		return err
	}

	u.Id = oid
	u.Username = user.Username
	u.Email = user.Email
	u.PhoneNumber = user.PhoneNumber
	return nil
}
