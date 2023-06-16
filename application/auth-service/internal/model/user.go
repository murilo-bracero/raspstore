package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id            primitive.ObjectID `bson:"_id"`
	UserId        string             `bson:"user_id"`
	Username      string             `bson:"username"`
	Password      string             `bson:"password"`
	Secret        string             `bson:"secret"`
	Permissions   []string           `bson:"permissions"`
	RefreshToken  string             `bson:"refresh_token"`
	IsMfaEnabled  bool               `bson:"is_mfa_enabled"`
	IsMfaVerified bool               `bson:"is_mfa_verified"`
}
