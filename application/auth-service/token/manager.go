package token

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"raspstore.github.io/auth-service/db"
)

type TokenManager interface {
	Verify(rawToken string) (uid string, err error)
	Generate(uid string) (token string, err error)
}

type tokenManager struct{}

func NewTokenManager() TokenManager {
	return &tokenManager{}
}

func (t *tokenManager) Verify(rawToken string) (string, error) {

	parsedToken, err := jwt.ParseWithClaims(rawToken, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("error reading jwt: wrong signing method: %v", token.Header["alg"])
		} else {
			return []byte(db.TokenSecret()), nil
		}
	})

	if err != nil {
		return "", err
	}

	return parsedToken.Claims.(*UserClaims).Uid, nil
}

func (t *tokenManager) Generate(uid string) (string, error) {

	claims := UserClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(db.TokenDuration()) * time.Second).Unix(),
		},
		Uid: uid,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(db.TokenSecret()))
}
