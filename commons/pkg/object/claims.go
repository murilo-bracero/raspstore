package object

import (
	"strings"

	"aidanwoods.dev/go-paseto"
)

type Claims struct {
	UserId string
	Roles  []string
}

func NewClaimsFromToken(token *paseto.Token) *Claims {
	decryptedMap := token.Claims()

	return &Claims{
		UserId: decryptedMap["sub"].(string),
		Roles:  strings.Split(decryptedMap["roles"].(string), ","),
	}
}
