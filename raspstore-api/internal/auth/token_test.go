package auth_test

import (
	"context"
	"testing"

	"github.com/murilo-bracero/raspstore/file-service/internal/auth"
	"github.com/murilo-bracero/raspstore/file-service/internal/infra/bootstrap"
	"github.com/murilo-bracero/raspstore/file-service/internal/infra/config"
	"github.com/stretchr/testify/assert"
)

func TestGenerateRefreshToken(t *testing.T) {

	mockConfig := &config.Config{Storage: config.StorageConfig{Path: t.TempDir()}}

	err := (&bootstrap.FolderBootstraper{}).Bootstrap(context.Background(), mockConfig)

	assert.NoError(t, err, "FolderBootstraper")

	err = (&bootstrap.SecretsBootstraper{}).Bootstrap(context.Background(), mockConfig)

	assert.NoError(t, err, "SecretsBootstraper")

	t.Run("should generate refresh token", func(t *testing.T) {
		token, err := auth.GenerateRefreshToken(mockConfig)

		assert.NoError(t, err, "GenerateRefreshToken")

		assert.NotEmpty(t, token, "GenerateRefreshToken")
	})
}
