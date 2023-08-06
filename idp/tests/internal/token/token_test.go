package token_test

import (
	"testing"

	"aidanwoods.dev/go-paseto"
	"github.com/google/uuid"
	"github.com/murilo-bracero/raspstore/idp/internal/infra"
	"github.com/murilo-bracero/raspstore/idp/internal/model"
	"github.com/murilo-bracero/raspstore/idp/internal/token"
	"github.com/stretchr/testify/assert"
)

var config *infra.Config

func init() {
	k := paseto.NewV4AsymmetricSecretKey()

	config = &infra.Config{
		TokenDuration:   12500,
		TokenPrivateKey: k.ExportHex(),
	}
}

func TestGenerateToken(t *testing.T) {
	user := &model.User{
		UserId:      uuid.NewString(),
		Permissions: []string{"admin"},
	}

	jwt, err := token.Generate(config, user)

	assert.NoError(t, err)

	assert.NotEmpty(t, jwt)
}
