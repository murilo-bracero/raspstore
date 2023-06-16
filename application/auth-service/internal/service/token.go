package service

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"raspstore.github.io/auth-service/internal"
	"raspstore.github.io/auth-service/internal/model"
)

type TokenService interface {
	Verify(rawToken string) (userClaims *model.UserClaims, err error)
	Generate(user *model.User) (string, error)
}

type tokenService struct{}

func NewTokenService() TokenService {
	return &tokenService{}
}

func (t *tokenService) Verify(rawToken string) (userClaims *model.UserClaims, err error) {

	parsedToken, err := jwt.ParseWithClaims(rawToken, &model.UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("error reading jwt: wrong signing method: %v", token.Header["alg"])
		} else {
			return []byte(internal.TokenSecret()), nil
		}
	})

	if err != nil {
		return nil, err
	}

	return parsedToken.Claims.(*model.UserClaims), nil
}

func (t *tokenService) Generate(user *model.User) (string, error) {

	claims := model.UserClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(internal.TokenDuration()) * time.Second).Unix(),
		},
		Uid:   user.UserId,
		Roles: user.Permissions,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(internal.TokenSecret()))
}
