package usecase_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/murilo-bracero/raspstore/file-service/internal/application/usecase"
	"github.com/murilo-bracero/raspstore/file-service/internal/domain/entity"
	"github.com/murilo-bracero/raspstore/file-service/internal/infra/repository"
	"github.com/murilo-bracero/raspstore/file-service/internal/infra/repository/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUpdateFileUseCase_Execute(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	_, err := createFile(uuid.NewString())
	assert.NoError(t, err)

	traceId := "trace12345"
	userId := "userId"

	t.Run("ValidFileUpdate", func(t *testing.T) {
		mockObj := mocks.NewMockTxFilesRepository(mockCtrl)

		mockObj.EXPECT().FindById(gomock.Any(), "userId", "validFile").Return(&entity.File{
			FileId:   "validFile",
			Filename: "updated.txt",
		}, nil)
		mockObj.EXPECT().Update(gomock.Any(), "userId", gomock.Any()).Return(nil)
		mockObj.EXPECT().Begin().Return(nil, nil)
		mockObj.EXPECT().Commit(nil).Return(nil)

		useCase := usecase.NewUpdateFileUseCase(mockObj)

		file := &entity.File{
			FileId:   "validFile",
			Filename: "updated.txt",
		}

		fileMetadata, err := useCase.Execute(file, userId, traceId)

		assert.NoError(t, err)
		assert.NotNil(t, fileMetadata)
		assert.Equal(t, "updated.txt", fileMetadata.Filename)
	})

	t.Run("FileNotFound", func(t *testing.T) {
		mockObj := mocks.NewMockTxFilesRepository(mockCtrl)

		mockObj.EXPECT().Begin().Return(nil, nil)
		mockObj.EXPECT().FindById(gomock.Any(), "userId", "nonexistentFile").Return(nil, repository.ErrFileDoesNotExists)

		useCase := usecase.NewUpdateFileUseCase(mockObj)

		file := &entity.File{
			FileId:   "nonexistentFile",
			Filename: "updated.txt",
		}

		fileMetadata, err := useCase.Execute(file, userId, traceId)

		assert.Error(t, err)
		assert.Nil(t, fileMetadata)
	})

	t.Run("FailedToUpdateFile", func(t *testing.T) {
		mockObj := mocks.NewMockTxFilesRepository(mockCtrl)

		mockObj.EXPECT().FindById(gomock.Any(), "userId", "failedFile").Return(&entity.File{
			FileId:   "failedFile",
			Filename: "updated.txt",
		}, nil)
		mockObj.EXPECT().Begin().Return(nil, nil)
		mockObj.EXPECT().Update(gomock.Any(), "userId", gomock.Any()).Return(errors.New("generic error"))

		useCase := usecase.NewUpdateFileUseCase(mockObj)

		file := &entity.File{
			FileId:   "failedFile",
			Filename: "updated.txt",
		}

		fileMetadata, err := useCase.Execute(file, userId, traceId)

		assert.Error(t, err)
		assert.Nil(t, fileMetadata)
	})
}

func createFile(seed string) (string, error) {
	dir := os.TempDir() + "/storage/"

	if err := os.Mkdir(dir, os.ModePerm); err != nil {
		if !os.IsExist(err) {
			return "", err
		}
	}

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
