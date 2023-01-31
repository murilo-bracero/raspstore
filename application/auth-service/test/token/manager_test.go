package token_test

import (
	"log"
	"testing"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"raspstore.github.io/auth-service/token"
)

func init() {
	err := godotenv.Load("../../.env")

	if err != nil {
		log.Panicln(err.Error())
	}
}

func TestGenerateToken(t *testing.T) {
	mngr := token.NewTokenManager()

	uid := uuid.NewString()

	token, err := mngr.Generate(uid)

	assert.NoError(t, err)

	assert.NotEmpty(t, token)

	if verified, err := mngr.Verify(token); err != nil {
		assert.Fail(t, err.Error())
	} else {
		assert.Equal(t, uid, verified)
	}
}

func TestFakeToken(t *testing.T) {
	mngr := token.NewTokenManager()

	token := "faketoken.token.fake"

	if _, err := mngr.Verify(token); err != nil {
		assert.Error(t, err)
	} else {
		assert.Fail(t, "accepted fraudulent token")
	}
}
