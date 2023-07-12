package model

import "github.com/dgrijalva/jwt-go"

type UserClaims struct {
	jwt.StandardClaims
	Uid   string   `json:"uid"`
	Roles []string `json:"roles"`
}
