package model

import (
	"time"

	"raspstore.github.io/authentication/pb"
)

type User struct {
	UserId      string    `bson:"user_id" datastore:"user_id"`
	Username    string    `bson:"username" datastore:"username"`
	Email       string    `bson:"email" datastore:"email"`
	PhoneNumber string    `bson:"phone_number" datastore:"phone_number"`
	CreatedAt   time.Time `bson:"created_at" datastore:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at" datastore:"updated_at"`
}

func (u *User) ToProtoBuffer() *pb.User {
	return &pb.User{
		Id:          u.UserId,
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
	u.UserId = user.Id
	u.Username = user.Username
	u.Email = user.Email
	u.PhoneNumber = user.PhoneNumber
	return nil
}
