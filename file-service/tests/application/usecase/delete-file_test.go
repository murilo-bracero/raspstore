package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/murilo-bracero/raspstore/file-service/internal/application/repository/mocks"
	"github.com/murilo-bracero/raspstore/file-service/internal/application/usecase"
	m "github.com/murilo-bracero/raspstore/file-service/internal/infra/middleware"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestDeleteFileUseCase(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	token := jwt.New()
	token.Set("sub", "userId")

	ctx := context.WithValue(context.Background(), middleware.RequestIDKey, "trace12345")
	ctx = context.WithValue(ctx, m.UserClaimsCtxKey, token)

	t.Run("ValidFileDeletion", func(t *testing.T) {
		mockObj := mocks.NewMockFilesRepository(mockCtrl)

		fileId := "file123"

		mockObj.EXPECT().Delete("userId", fileId).Return(nil)

		useCase := usecase.NewDeleteFileUseCase(mockObj)

		err := useCase.Execute(ctx, fileId)

		assert.NoError(t, err)
	})

	t.Run("FileDeletionError", func(t *testing.T) {
		mockObj := mocks.NewMockFilesRepository(mockCtrl)

		fileId := "file123"

		mockObj.EXPECT().Delete("userId", fileId).Return(errors.New("generic error"))

		useCase := usecase.NewDeleteFileUseCase(mockObj)

		err := useCase.Execute(ctx, fileId)

		assert.Error(t, err)
	})
}
