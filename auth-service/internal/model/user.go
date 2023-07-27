package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id            primitive.ObjectID `bson:"_id"`
	UserId        string             `bson:"user_id"`
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
