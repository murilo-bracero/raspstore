package model

import (
	"time"
)

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
