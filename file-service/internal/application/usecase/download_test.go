package usecase_test

import (
	"context"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"github.com/murilo-bracero/raspstore/file-service/internal/application/usecase"
	"github.com/stretchr/testify/assert"
)

func createFile(seed string) (string, error) {
	dir := os.TempDir() + "/storage/"

	if err := os.Mkdir(dir, os.ModePerm); err != nil {
		slog.Error("could not create storage temp folder", "err", err)
	}

	slog.Info("creating temp file", "dir", dir, "seed", seed)

	file, err := os.CreateTemp(dir, seed)

	if err != nil {
		return "", err
	}

	defer file.Close()

	if _, err := file.WriteString("test file content"); err != nil {
		return "", err
	}

	return filepath.Base(file.Name()), nil
}

func TestDownloadFileUseCase(t *testing.T) {
	mockConfig.Storage.Path = os.TempDir()

	seed := uuid.NewString()
	fileId, err := createFile(seed)
	assert.NoError(t, err)

	ctx := context.Background()
	ctx = context.WithValue(ctx, middleware.RequestIDKey, "test-trace-id")

	t.Cleanup(func() {
		if err := os.RemoveAll(os.TempDir() + "/storage"); err != nil {
			slog.Error("could not cleanup temp folder", "err", err)
		}
	})

	t.Run("happy path", func(t *testing.T) {
		uc := usecase.NewDownloadFileUseCase(mockConfig)

		file, err := uc.Execute(ctx, fileId)

		assert.NoError(t, err)

		assert.NotNil(t, file)

		content, err := io.ReadAll(file)

		assert.NoError(t, err)

		assert.NotNil(t, file)

		assert.Equal(t, "test file content", string(content))
	})

	t.Run("should return error when file does not exists", func(t *testing.T) {
		uc := usecase.NewDownloadFileUseCase(mockConfig)

		_, err := uc.Execute(ctx, "no-exists")

		assert.Error(t, err)
	})
}
