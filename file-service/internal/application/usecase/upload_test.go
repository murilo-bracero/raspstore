package usecase_test

import (
	"context"
	"os"
	"testing"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"github.com/murilo-bracero/raspstore/file-service/internal/application/usecase"
	"github.com/murilo-bracero/raspstore/file-service/internal/domain/entity"
	"github.com/stretchr/testify/assert"
)

func TestUploadFileUseCase(t *testing.T) {
	mockConfig.Server.Storage.Path = "/tmp"

	eFile := &entity.File{
		FileId: uuid.NewString(),
	}

	file, err := os.Create("/tmp/storage/upload.txt")

	assert.NoError(t, err)

	ctx := context.Background()
	ctx = context.WithValue(ctx, middleware.RequestIDKey, "test-trace-id")

	t.Run("happy path", func(t *testing.T) {
		uc := usecase.NewUploadFileUseCase(mockConfig)

		err := uc.Execute(ctx, eFile, file)

		assert.NoError(t, err)

		_, err = os.Lstat(mockConfig.Server.Storage.Path + "/storage/" + eFile.FileId)

		assert.NoError(t, err)
	})
}
