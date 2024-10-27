package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHash(t *testing.T) {

	hash, err := HashPassword("test")

	assert.NoError(t, err, "HashPassword")

	assert.NotEmpty(t, hash, "HashPassword")
}
