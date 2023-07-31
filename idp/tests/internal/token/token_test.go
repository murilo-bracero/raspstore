package token_test

import (
	"log"
	"testing"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/murilo-bracero/raspstore/idp/internal/infra"
	"github.com/murilo-bracero/raspstore/idp/internal/model"
	"github.com/murilo-bracero/raspstore/idp/internal/token"
	"github.com/stretchr/testify/assert"
)

var config *infra.Config

func init() {
	err := godotenv.Load("../../.env.test")

	if err != nil {
		log.Panicln(err.Error())
	}

	config = infra.NewConfig()
}

func TestGenerateToken(t *testing.T) {
	user := &model.User{
		UserId:      uuid.NewString(),
		Permissions: []string{"admin"},
	}

	jwt, err := token.Generate(config, user)

	assert.NoError(t, err)

	assert.NotEmpty(t, jwt)

	if claims, err := token.Verify(config, jwt); err != nil {
		assert.Fail(t, err.Error())
	} else {
		assert.Equal(t, user.Permissions, claims.Roles)
		assert.Equal(t, user.UserId, claims.Uid)
	}
}

func TestFakeToken(t *testing.T) {
	jwt := "faketoken.token.fake"

	if _, err := token.Verify(config, jwt); err != nil {
		assert.Error(t, err)
	} else {
		assert.Fail(t, "accepted fraudulent token")
	}
}
