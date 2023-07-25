package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/murilo-bracero/raspstore/auth-service/internal"
	"github.com/murilo-bracero/raspstore/auth-service/internal/model"
)

func Generate(user *model.User) (string, error) {
	claims := model.UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.NewString(),
			IssuedAt:  &jwt.NumericDate{Time: time.Now()},
			Audience:  jwt.ClaimStrings{"account"},
			Subject:   user.UserId,
			NotBefore: &jwt.NumericDate{Time: time.Now()},
			ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(time.Duration(internal.TokenDuration()) * time.Second)},
		},
		Roles: user.Permissions,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(internal.TokenSecret()))
}

func Verify(token string) (userClaims *model.UserClaims, err error) {
	parsedToken, err := jwt.ParseWithClaims(token, &model.UserClaims{}, func(jt *jwt.Token) (interface{}, error) {
		if _, ok := jt.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("error reading jwt: wrong signing method: %v", jt.Header["alg"])
		}
		return []byte(internal.TokenSecret()), nil
	})

	if err != nil {
		return nil, err
	}

	return parsedToken.Claims.(*model.UserClaims), nil
}
