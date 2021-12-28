package model

import (
	"time"

	api "raspstore.github.io/authentication/api/dto"
	"raspstore.github.io/authentication/pb"
)

type User struct {
	UserId      string    `bson:"user_id" json:"userId"`
	Username    string    `bson:"username" json:"username"`
	Email       string    `bson:"email" json:"email"`
	PhoneNumber string    `bson:"phone_number" json:"phoneNumber"`
	CreatedAt   time.Time `bson:"created_at" json:"createdAt"`
	UpdatedAt   time.Time `bson:"updated_at" json:"updatedAt"`
}

func (u *User) FromUpdateRequest(uur api.UpdateUserRequest, id string) {
	u.UserId = id
	u.Username = uur.Username
	u.Email = uur.Email
	u.PhoneNumber = uur.PhoneNumber
}

func (u *User) FromCreateRequest(cur api.CreateUserRequest) {

	u.Username = cur.Username
	u.Email = cur.Email
	u.PhoneNumber = cur.PhoneNumber
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
