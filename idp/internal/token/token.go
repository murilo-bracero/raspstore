package token

import (
	"strings"
	"time"

	"aidanwoods.dev/go-paseto"
	"github.com/google/uuid"
	"github.com/murilo-bracero/raspstore/idp/internal/infra"
	"github.com/murilo-bracero/raspstore/idp/internal/model"
)

func Generate(config *infra.Config, user *model.User) (string, error) {
	token := paseto.NewToken()

	token.SetJti(uuid.NewString())
	token.SetIssuedAt(time.Now())
	token.SetNotBefore(time.Now())
	token.SetExpiration(time.Now().Add((time.Duration(config.TokenDuration) * time.Hour)))
	token.SetSubject(user.UserId)
	token.SetAudience("account")
	token.SetString("roles", strings.Join(user.Permissions, ","))

	key, err := paseto.V4SymmetricKeyFromHex(config.TokenSecret)

	if err != nil {
		return "", err
	}

	return token.V4Encrypt(key, nil), nil
}

func Verify(config *infra.Config, token string) (claims *model.UserClaims, err error) {
	parser := paseto.NewParser()

	parser.AddRule(paseto.ForAudience("account"))
	parser.AddRule(paseto.NotExpired())
	parser.AddRule(paseto.ValidAt(time.Now()))

	key, err := paseto.V4SymmetricKeyFromHex(config.TokenSecret)

	if err != nil {
		return nil, err
	}

	decrypted, err := parser.ParseV4Local(key, token, nil)

	if err != nil {
		return nil, err
	}

	decryptedMap := decrypted.Claims()

	return &model.UserClaims{
		Uid:   decryptedMap["sub"].(string),
		Roles: strings.Split(decryptedMap["roles"].(string), ","),
	}, nil
}
