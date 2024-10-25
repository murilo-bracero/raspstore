package facade

import (
	"io"
	"log/slog"
	"os"

	"github.com/murilo-bracero/raspstore/file-service/internal/domain/entity"
	"github.com/murilo-bracero/raspstore/file-service/internal/infra/config"
)

type FileSystemFacade interface {
	Download(traceId string, requesterId string, fileId string) (*os.File, error)
	Upload(traceId string, file *entity.File, src io.Reader) error
}

type fileSystemFacade struct {
	config *config.Config
}

var _ FileSystemFacade = (*fileSystemFacade)(nil)

func NewFileSystemFacade(c *config.Config) *fileSystemFacade {
	return &fileSystemFacade{config: c}
}

func (f *fileSystemFacade) Download(traceId string, requesterId string, fileId string) (*os.File, error) {
	file, err := os.Open(f.config.Storage.Path + "/storage/" + fileId)

	if err != nil {
		slog.Error("Could not open file in fs", "traceId", traceId, "error", err)
		return nil, err
	}

	return file, nil
}

func (f *fileSystemFacade) Upload(traceId string, file *entity.File, src io.Reader) error {
	filerep, err := os.Create(f.config.Storage.Path + "/storage/" + file.FileId)

	if err != nil {
		slog.Error("Could not create file in fs", "traceId", traceId, "error", err)
		return err
	}

	defer filerep.Close()

	_, err = io.Copy(filerep, src)

	if err != nil {
		slog.Error("Could not read file buffer", "traceId", traceId, "error", err)
		return err
	}

	return nil
}
