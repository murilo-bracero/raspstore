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

	key, err := paseto.NewV4AsymmetricSecretKeyFromHex(config.TokenPrivateKey)

	if err != nil {
		return "", err
	}

	return token.V4Sign(key, nil), nil
}
