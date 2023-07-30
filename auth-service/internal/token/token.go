package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/murilo-bracero/raspstore/auth-service/internal/infra"
	"github.com/murilo-bracero/raspstore/auth-service/internal/model"
)

func Generate(config *infra.Config, user *model.User) (string, error) {
	claims := model.UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.NewString(),
			IssuedAt:  &jwt.NumericDate{Time: time.Now()},
			Audience:  jwt.ClaimStrings{"account"},
			Subject:   user.UserId,
			NotBefore: &jwt.NumericDate{Time: time.Now()},
			ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(time.Duration(config.TokenDuration) * time.Second)},
		},
		Roles: user.Permissions,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(config.TokenSecret))
}

func Verify(config *infra.Config, token string) (userClaims *model.UserClaims, err error) {
	parsedToken, err := jwt.ParseWithClaims(token, &model.UserClaims{}, func(jt *jwt.Token) (interface{}, error) {
		if _, ok := jt.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("error reading jwt: wrong signing method: %v", jt.Header["alg"])
		}
		return []byte(config.TokenSecret), nil
	})

	if err != nil {
		return nil, err
	}

	return parsedToken.Claims.(*model.UserClaims), nil
}
