package service

import (
	"log"
	"testing"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"raspstore.github.io/auth-service/internal/service"
)

func init() {
	err := godotenv.Load("../../.env.test")

	if err != nil {
		log.Panicln(err.Error())
	}
}

func TestGenerateToken(t *testing.T) {
	ts := service.NewTokenService()

	uid := uuid.NewString()

	token, err := ts.Generate(uid)

	assert.NoError(t, err)

	assert.NotEmpty(t, token)

	if verified, err := ts.Verify(token); err != nil {
		assert.Fail(t, err.Error())
	} else {
		assert.Equal(t, uid, verified)
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
