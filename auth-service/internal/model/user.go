package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	UserId        string `bson:"user_id"`
	Username      string
	IsEnabled     bool `bson:"is_enabled"`
	Password      string
	Secret        string
	Permissions   []string
	RefreshToken  string    `bson:"refresh_token"`
	IsMfaEnabled  bool      `bson:"is_mfa_enabled"`
	IsMfaVerified bool      `bson:"is_mfa_verified"`
	CreatedAt     time.Time `bson:"created_at"`
	UpdatedAt     time.Time `bson:"updated_at"`
}

func (u *User) BeforeCreate() {
	u.UserId = uuid.NewString()
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
}

type UserPage struct {
	Content []*User `bson:"content"`
	Count   int     `bson:"count"`
}
