package usecase_test

import (
	"context"
	"io"
	"os"
	"testing"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"github.com/murilo-bracero/raspstore/file-service/internal/application/usecase"
	"github.com/stretchr/testify/assert"
)

func createFile(fileId string) error {
	file, err := os.Create("/tmp/" + fileId)

	if err != nil {
		return err
	}

	defer file.Close()

	if _, err := file.WriteString("test file content"); err != nil {
		return err
	}

	return nil
}

func TestDownloadFileUseCase(t *testing.T) {
	mockConfig.Server.Storage.Path = "/tmp"

	fileId := uuid.New()
	err := createFile(fileId.String())
	assert.NoError(t, err)

	ctx := context.Background()
	ctx = context.WithValue(ctx, middleware.RequestIDKey, "test-trace-id")

	t.Run("happy path", func(t *testing.T) {
		uc := usecase.NewDownloadFileUseCase(mockConfig)

		file, err := uc.Execute(ctx, fileId.String())

		assert.NoError(t, err)

		assert.NotNil(t, file)

		content, err := io.ReadAll(file)

		assert.NotNil(t, file)

		assert.Equal(t, "test file content", string(content))
	})

	t.Run("should return error when file does not exists", func(t *testing.T) {
		uc := usecase.NewDownloadFileUseCase(mockConfig)

		_, err := uc.Execute(ctx, "no-exists")

		assert.Error(t, err)
	})
}
