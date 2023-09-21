package usecase_test

import (
	"context"
	"testing"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/murilo-bracero/raspstore/file-service/internal/application/usecase"
	m "github.com/murilo-bracero/raspstore/file-service/internal/infra/middleware"
	"github.com/stretchr/testify/assert"
)

func TestDeleteFileUseCase(t *testing.T) {
	useCase := usecase.NewDeleteFileUseCase(&mockFilesRepository{})

	token := jwt.New()
	token.Set("sub", uuid.NewString())

	ctx := context.WithValue(context.WithValue(context.Background(),
		middleware.RequestIDKey, "trace12345"),
		m.UserClaimsCtxKey, token)

	t.Run("ValidFileDeletion", func(t *testing.T) {
		fileId := "file123"

		err := useCase.Execute(ctx, fileId)

		assert.NoError(t, err)
	})

	t.Run("FileDeletionError", func(t *testing.T) {
		useCase := usecase.NewDeleteFileUseCase(&mockFilesRepository{})

		fileId := "file456_error"

		err := useCase.Execute(ctx, fileId)

		assert.Error(t, err)
	})
}
