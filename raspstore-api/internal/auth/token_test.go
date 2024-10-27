package auth

import (
	"context"
	"os"
	"testing"

	"github.com/murilo-bracero/raspstore/file-service/internal/infra/bootstrap"
	"github.com/murilo-bracero/raspstore/file-service/internal/infra/config"
	"github.com/stretchr/testify/assert"
)

func TestGenerateRefreshToken(t *testing.T) {

	mockConfig := &config.Config{Storage: config.StorageConfig{Path: os.TempDir()}}

	(&bootstrap.FolderBootstraper{}).Bootstrap(context.Background(), mockConfig)

	(&bootstrap.SecretsBootstraper{}).Bootstrap(context.Background(), mockConfig)

	t.Cleanup(func() {
		os.RemoveAll(os.TempDir() + "/secrets")
	})

	t.Run("should generate refresh token", func(t *testing.T) {
		token, err := GenerateRefreshToken(mockConfig)

		assert.NoError(t, err, "GenerateRefreshToken")

		assert.NotEmpty(t, token, "GenerateRefreshToken")
	})
}
