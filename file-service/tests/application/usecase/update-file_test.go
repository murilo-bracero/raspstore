package usecase_test

import (
	"context"
	"errors"
	"testing"

	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/murilo-bracero/raspstore/file-service/internal/application/repository/mocks"
	"github.com/murilo-bracero/raspstore/file-service/internal/application/usecase"
	"github.com/murilo-bracero/raspstore/file-service/internal/domain/entity"
	"github.com/murilo-bracero/raspstore/file-service/internal/domain/exceptions"
	"github.com/murilo-bracero/raspstore/file-service/internal/infra/middleware"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUpdateFileUseCase_Execute(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	token := jwt.New()
	token.Set("sub", "userId")

	ctx := context.WithValue(context.WithValue(context.Background(),
		chiMiddleware.RequestIDKey, "trace12345"),
		middleware.UserClaimsCtxKey, token)

	t.Run("ValidFileUpdate", func(t *testing.T) {
		mockObj := mocks.NewMockFilesRepository(mockCtrl)

		mockObj.EXPECT().FindById("userId", "validFile").Return(&entity.File{
			FileId:   "validFile",
			Filename: "updated.txt",
		}, nil)
		mockObj.EXPECT().Update("userId", gomock.Any()).Return(nil)

		useCase := usecase.NewUpdateFileUseCase(mockObj)

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
		mockObj := mocks.NewMockFilesRepository(mockCtrl)

		mockObj.EXPECT().FindById("userId", "nonexistentFile").Return(nil, exceptions.ErrFileDoesNotExists)

		useCase := usecase.NewUpdateFileUseCase(mockObj)

		file := &entity.File{
			FileId:   "nonexistentFile",
			Filename: "updated.txt",
		}

		fileMetadata, err := useCase.Execute(ctx, file)

		assert.Error(t, err)
		assert.Nil(t, fileMetadata)
	})

	t.Run("FailedToUpdateFile", func(t *testing.T) {
		mockObj := mocks.NewMockFilesRepository(mockCtrl)

		mockObj.EXPECT().FindById("userId", "failedFile").Return(&entity.File{
			FileId:   "failedFile",
			Filename: "updated.txt",
		}, nil)
		mockObj.EXPECT().Update("userId", gomock.Any()).Return(errors.New("generic error"))

		useCase := usecase.NewUpdateFileUseCase(mockObj)

		file := &entity.File{
			FileId:   "failedFile",
			Filename: "updated.txt",
		}

		fileMetadata, err := useCase.Execute(ctx, file)

		assert.Error(t, err)
		assert.Nil(t, fileMetadata)
	})
}
