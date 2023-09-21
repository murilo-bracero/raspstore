package usecase_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/murilo-bracero/raspstore/file-service/internal/application/usecase"
	"github.com/murilo-bracero/raspstore/file-service/internal/domain/entity"
	"github.com/murilo-bracero/raspstore/file-service/internal/infra/config"
)

type mockFilesRepository struct{}

func (m *mockFilesRepository) FindUsageByUserId(userID string) (int64, error) {
	return 100, nil
}

func (m *mockFilesRepository) Save(file *entity.File) error {
	return nil
}

func (m *mockFilesRepository) Delete(userId string, fileId string) error {
	if strings.HasSuffix(fileId, "_error") {
		return errors.New("generic error")
	}
	return nil
}

func (m *mockFilesRepository) FindAll(userId string, page int, size int, filename string, secret bool) (filesPage *entity.FilePage, err error) {
	return nil, errors.New("Not implemented!")
}

func (m *mockFilesRepository) FindById(userId string, fileId string) (*entity.File, error) {
	if fileId == "validFile" {
		return &entity.File{
			FileId:   "validFile",
			Filename: "example.txt",
		}, nil
	}
	return nil, errors.New("file not found")
}

func (m *mockFilesRepository) Update(userId string, file *entity.File) error {
	if file.FileId == "validFile" {
		return nil
	}
	return errors.New("failed to update file")
}

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

	useCase := usecase.NewCreateFileUseCase(mockConfig, &mockFilesRepository{})

	t.Run("happy path", func(t *testing.T) {
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
