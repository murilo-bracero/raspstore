package usecase_test

import (
	"context"
	"testing"

	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/murilo-bracero/raspstore/file-service/internal/application/usecase"
	"github.com/murilo-bracero/raspstore/file-service/internal/domain/entity"
	"github.com/murilo-bracero/raspstore/file-service/internal/infra/middleware"
	"github.com/stretchr/testify/assert"
)

func TestUpdateFileUseCase_Execute(t *testing.T) {
	useCase := usecase.NewUpdateFileUseCase(&mockFilesRepository{})

	token := jwt.New()
	token.Set("sub", uuid.NewString())

	ctx := context.WithValue(context.WithValue(context.Background(),
		chiMiddleware.RequestIDKey, "trace12345"),
		middleware.UserClaimsCtxKey, token)

	t.Run("ValidFileUpdate", func(t *testing.T) {
		file := &entity.File{
			FileId:   "validFile",
			Filename: "updated.txt",
		}

		fileMetadata, err := useCase.Execute(ctx, file)

		assert.NoError(t, err)
		assert.NotNil(t, fileMetadata)
		assert.Equal(t, "updated.txt", fileMetadata.Filename)
	})

	t.Run("FileNotFound", func(t *testing.T) {
		file := &entity.File{
			FileId:   "nonexistentFile",
			Filename: "updated.txt",
		}

		fileMetadata, err := useCase.Execute(ctx, file)

		assert.Error(t, err)
		assert.Nil(t, fileMetadata)
	})

	t.Run("FailedToUpdateFile", func(t *testing.T) {
		useCase := usecase.NewUpdateFileUseCase(&mockFilesRepository{})

		file := &entity.File{
			FileId:   "failedFile",
			Filename: "updated.txt",
		}

		fileMetadata, err := useCase.Execute(ctx, file)

		assert.Error(t, err)
		assert.Nil(t, fileMetadata)
	})
}
