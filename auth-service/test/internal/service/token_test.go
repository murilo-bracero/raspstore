package service

import (
	"log"
	"testing"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/murilo-bracero/raspstore/auth-service/internal/model"
	"github.com/murilo-bracero/raspstore/auth-service/internal/service"
	"github.com/stretchr/testify/assert"
)

func init() {
	err := godotenv.Load("../../.env.test")

	if err != nil {
		log.Panicln(err.Error())
	}
}

func TestGenerateToken(t *testing.T) {
	ts := service.NewTokenService()

	user := &model.User{
		UserId:      uuid.NewString(),
		Permissions: []string{"admin"},
	}

	token, err := ts.Generate(user)

	assert.NoError(t, err)

	assert.NotEmpty(t, token)

	if claims, err := ts.Verify(token); err != nil {
		assert.Fail(t, err.Error())
	} else {
		assert.Equal(t, user.Permissions, claims.Roles)
		assert.Equal(t, user.UserId, claims.Uid)
	}
}

func TestFakeToken(t *testing.T) {
	ts := service.NewTokenService()

	token := "faketoken.token.fake"

	if _, err := ts.Verify(token); err != nil {
		assert.Error(t, err)
	} else {
		assert.Fail(t, "accepted fraudulent token")
	}
}
