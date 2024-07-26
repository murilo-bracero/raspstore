package usecase_test

import (
	"testing"

	"github.com/murilo-bracero/raspstore/file-service/internal/application/repository/mocks"
	"github.com/murilo-bracero/raspstore/file-service/internal/application/usecase"
	"github.com/murilo-bracero/raspstore/file-service/internal/domain/entity"
	"github.com/murilo-bracero/raspstore/file-service/internal/infra/config"
	"go.uber.org/mock/gomock"
)

var mockConfig = &config.Config{
	Server: struct {
		Port    int
		Storage struct {
			Path  string
			Limit string
		}
	}{Storage: struct {
		Path  string
		Limit string
	}{Path: "./", Limit: "1000M"}}}

func TestCreateFileUseCase(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	t.Run("happy path", func(t *testing.T) {
		mockObj := mocks.NewMockFilesRepository(mockCtrl)
		mockObj.EXPECT().Save(gomock.Any()).Return(nil)
		mockObj.EXPECT().FindUsageByUserId("user1").Return(int64(100), nil)

		useCase := usecase.NewCreateFileUseCase(mockConfig, mockObj)

		file := &entity.File{
			Owner: "user1",
			Size:  toMb(100),
		}

		err := useCase.Execute(file)

		if err != nil {
			t.Errorf("Expected no error, but got: %v", err)
		}
	})

	t.Run("upload with file size greather than provided by config", func(t *testing.T) {
		mockObj := mocks.NewMockFilesRepository(mockCtrl)
		mockObj.EXPECT().FindUsageByUserId("user2").Return(int64(100), nil)

		useCase := usecase.NewCreateFileUseCase(mockConfig, mockObj)

		file := &entity.File{
			Owner: "user2",
			Size:  toMb(1001),
		}

		err := useCase.Execute(file)

		if err != usecase.ErrNotAvailableSpace {
			t.Errorf("Expected ErrNotAvailableSpace, but got: %v", err)
		}
	})
}

func toMb(v int64) int64 {
	return v * 1024 * 1024
}
