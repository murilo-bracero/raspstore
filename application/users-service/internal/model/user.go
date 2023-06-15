package model

import (
	"time"

	"github.com/google/uuid"
	v1 "raspstore.github.io/users-service/api/v1"
)

const defaultDateFormat = "2006-01-02 15:04:05"

type User struct {
	UserId       string    `bson:"user_id"`
	Username     string    `bson:"username"`
	Email        string    `bson:"email"`
	PasswordHash string    `bson:"password"`
	IsEnabled    bool      `bson:"is_enabled"`
	PhoneNumber  string    `bson:"phone_number"`
	Permissions  []string  `bson:"permissions"`
	CreatedAt    time.Time `bson:"created_at"`
	UpdatedAt    time.Time `bson:"updated_at"`
}

func NewUserByCreateUserRequest(req v1.CreateUserRequest) *User {
	return &User{
		UserId:       uuid.NewString(),
		Username:     req.Username,
		Email:        req.Email,
		IsEnabled:    true,
		PhoneNumber:  "",
		Permissions:  []string{},
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		PasswordHash: req.Password,
	}
}

func NewUserByAdminCreateUserRequest(req v1.AdminCreateUserRequest) *User {
	return &User{
		UserId:       uuid.NewString(),
		Username:     req.Username,
		Email:        req.Email,
		IsEnabled:    true,
		PhoneNumber:  "",
		Permissions:  req.Permissions,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		PasswordHash: req.Password,
	}
}

func (usr *User) ToDto() v1.UserResponse {
	return v1.UserResponse{
		UserId:    usr.UserId,
		Username:  usr.Username,
		Email:     usr.Email,
		CreatedAt: usr.CreatedAt.Format(defaultDateFormat),
		UpdatedAt: usr.UpdatedAt.Format(defaultDateFormat),
	}
}

type UserPage struct {
	Content []*User `bson:"content"`
	Count   int     `bson:"count"`
}
