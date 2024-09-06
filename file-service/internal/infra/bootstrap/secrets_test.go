package bootstrap_test

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"testing"

	"github.com/murilo-bracero/raspstore/file-service/internal/infra/bootstrap"
	"github.com/murilo-bracero/raspstore/file-service/internal/infra/config"
	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	config := &config.Config{
		Storage: struct {
			Path  string
			Limit string
		}{Path: os.TempDir()},
	}

	err := os.Mkdir(os.TempDir()+"/secrets", os.ModePerm)

	assert.NoError(t, err, "os.Makedir")

	ctx := context.Background()

	t.Run("should create new keys if directory is empty", func(t *testing.T) {
		bt := &bootstrap.SecretsBootstrap{}

		bt.Bootstrap(ctx, config)

		filename := os.TempDir() + "/secrets/local-jwk.json"

		fi, err := os.Stat(filename)

		assert.NoError(t, err, "os.Stat")

		assert.NotEqual(t, 0, fi.Size())

		fl, err := os.Open(filename)

		assert.NoError(t, err, "os.Open")

		content, err := io.ReadAll(fl)

		var jsonKeyMap map[string]interface{}
		err = json.Unmarshal(content, &jsonKeyMap)

		assert.NoError(t, err)

		assert.Equal(t, "AQAB", jsonKeyMap["e"])
		assert.Equal(t, "RSA", jsonKeyMap["kty"])
		assert.NotEmpty(t, jsonKeyMap["d"])
		assert.NotEmpty(t, jsonKeyMap["dp"])
		assert.NotEmpty(t, jsonKeyMap["dq"])
		assert.NotEmpty(t, jsonKeyMap["n"])
		assert.NotEmpty(t, jsonKeyMap["p"])
		assert.NotEmpty(t, jsonKeyMap["q"])
		assert.NotEmpty(t, jsonKeyMap["qi"])
	})

	t.Cleanup(func() {
		os.RemoveAll(os.TempDir() + "/secrets")
	})
}
