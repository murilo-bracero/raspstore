package facade_test

import (
	"io"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/murilo-bracero/raspstore/file-service/internal/application/facade"
	"github.com/murilo-bracero/raspstore/file-service/internal/domain/entity"
	"github.com/murilo-bracero/raspstore/file-service/internal/infra/config"
	"github.com/stretchr/testify/assert"
)

func TestUploadFileUseCase(t *testing.T) {
	config := &config.Config{Storage: config.StorageConfig{Path: os.TempDir()}}

	eFile := &entity.File{
		FileId: uuid.NewString(),
	}

	storagePath := path.Join(os.TempDir(), "storage")

	err := os.MkdirAll(storagePath, os.ModePerm)

	if err != nil && !os.IsExist(err) {
		assert.Fail(t, "os.MkdirAll")
	}

	t.Cleanup(func() {
		err := os.RemoveAll(os.TempDir() + "/storage")

		assert.NoError(t, err, "os.RemoveAll")
	})

	file, err := os.CreateTemp(storagePath, "upload.*.txt")

	assert.NoError(t, err)

	t.Run("happy path", func(t *testing.T) {
		ffc := facade.NewFileSystemFacade(config)

		err := ffc.Upload("test-trace-id", eFile, file)

		assert.NoError(t, err)

		_, err = os.Lstat(config.Storage.Path + "/storage/" + eFile.FileId)

		assert.NoError(t, err)
	})
}

func TestDownloadFileUseCase(t *testing.T) {
	config := &config.Config{Storage: config.StorageConfig{Path: os.TempDir()}}

	seed := uuid.NewString()
	fileId, err := createFile(seed)
	assert.NoError(t, err)

	t.Cleanup(func() {
		err := os.RemoveAll(os.TempDir() + "/storage")

		assert.NoError(t, err, "os.RemoveAll")
	})

	t.Run("happy path", func(t *testing.T) {
		ffc := facade.NewFileSystemFacade(config)

		file, err := ffc.Download("test-trace-id", "userId", fileId)

		assert.NoError(t, err)

		assert.NotNil(t, file)

		content, err := io.ReadAll(file)

		assert.NoError(t, err)

		assert.NotNil(t, file)

		assert.Equal(t, "test file content", string(content))
	})

	t.Run("should return error when file does not exists", func(t *testing.T) {
		ffc := facade.NewFileSystemFacade(config)

		_, err := ffc.Download("test-trace-id", "userId", "no-exists")

		assert.Error(t, err)
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
