package model

import (
	"time"

	"github.com/google/uuid"
	"raspstore.github.io/users-service/api/dto"
)

const defaultDateFormat = "2006-01-02 15:04:05"

type User struct {
	UserId       string    `bson:"user_id"`
	Username     string    `bson:"username"`
	Email        string    `bson:"email"`
	PasswordHash string    `bson:"password"`
	IsEnabled    bool      `bson:"is_enabled"`
	PhoneNumber  string    `bson:"phone_number"`
	CreatedAt    time.Time `bson:"created_at"`
	UpdatedAt    time.Time `bson:"updated_at"`
}

func NewUserByCreateUserRequest(req dto.CreateUserRequest) *User {
	return &User{
		UserId:      uuid.NewString(),
		Username:    req.Username,
		Email:       req.Email,
		IsEnabled:   true,
		PhoneNumber: "",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func (usr *User) ToDto() dto.UserResponse {
	return dto.UserResponse{
		UserId:    usr.UserId,
		Username:  usr.Username,
		Email:     usr.Email,
		IsEnabled: usr.IsEnabled,
		CreatedAt: usr.CreatedAt.Format(defaultDateFormat),
		UpdatedAt: usr.UpdatedAt.Format(defaultDateFormat),
	}
}

type UserPage struct {
	Content []*User `bson:"content"`
	Count   int     `bson:"count"`
}
