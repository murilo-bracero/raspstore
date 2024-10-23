package facade_test

import (
	"testing"

	"github.com/murilo-bracero/raspstore/file-service/internal/application/facade"
	"github.com/murilo-bracero/raspstore/file-service/internal/domain/entity"
	"github.com/murilo-bracero/raspstore/file-service/internal/infra/config"
	"github.com/murilo-bracero/raspstore/file-service/internal/infra/repository/mocks"
	"go.uber.org/mock/gomock"
)

func TestSave(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	mockConfig := &config.Config{Storage: config.StorageConfig{Path: "./", Limit: "1000M"}}

	t.Run("happy path", func(t *testing.T) {
		mockObj := mocks.NewMockFilesRepository(mockCtrl)
		mockObj.EXPECT().Save(gomock.Any()).Return(nil)
		mockObj.EXPECT().FindUsageByUserId("user1").Return(int64(100), nil)

		ff := facade.NewFileFacade(mockConfig, mockObj)

		file := &entity.File{
			Owner: "user1",
			Size:  toMb(100),
		}

		err := ff.Save(file)

		if err != nil {
			t.Errorf("Expected no error, but got: %v", err)
		}
	})

	t.Run("upload with file size greather than provided by config", func(t *testing.T) {
		mockObj := mocks.NewMockFilesRepository(mockCtrl)
		mockObj.EXPECT().FindUsageByUserId("user2").Return(int64(100), nil)

		ff := facade.NewFileFacade(mockConfig, mockObj)

		file := &entity.File{
			Owner: "user2",
			Size:  toMb(1001),
		}

		err := ff.Save(file)

		if err != facade.ErrNotAvailableSpace {
			t.Errorf("Expected ErrNotAvailableSpace, but got: %v", err)
		}
	})
}

func toMb(v int64) int64 {
	return v * 1024 * 1024
}
