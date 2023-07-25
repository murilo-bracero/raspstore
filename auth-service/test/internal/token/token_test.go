package token_test

import (
	"log"
	"testing"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/murilo-bracero/raspstore/auth-service/internal/model"
	"github.com/murilo-bracero/raspstore/auth-service/internal/token"
	"github.com/stretchr/testify/assert"
)

func init() {
	err := godotenv.Load("../../.env.test")

	if err != nil {
		log.Panicln(err.Error())
	}
}

func TestGenerateToken(t *testing.T) {
	user := &model.User{
		UserId:      uuid.NewString(),
		Permissions: []string{"admin"},
	}

	jwt, err := token.Generate(user)

	assert.NoError(t, err)

	assert.NotEmpty(t, jwt)

	if claims, err := token.Verify(jwt); err != nil {
		assert.Fail(t, err.Error())
	} else {
		assert.Equal(t, user.Permissions, claims.Roles)
		assert.Equal(t, user.UserId, claims.Subject)
	}
}

func TestFakeToken(t *testing.T) {
	jwt := "faketoken.token.fake"

	if _, err := token.Verify(jwt); err != nil {
		assert.Error(t, err)
	} else {
		assert.Fail(t, "accepted fraudulent token")
	}
}
