package facade

import (
	"log/slog"

	"github.com/murilo-bracero/raspstore/file-service/internal/application/repository"
	"github.com/murilo-bracero/raspstore/file-service/internal/domain/entity"
)

const maxListSize = 50

type FileFacade interface {
	FindById(requesterId string, fileId string) (*entity.File, error)
	DeleteById(traceId string, requesterId string, fileId string) error
	FindAll(traceId string, requesterId string, page int, size int, filename string, secret bool) (*entity.FilePage, error)
}

type fileFacade struct {
	filesRepository repository.FilesRepository
}

func NewFileFacade(filesRepository repository.FilesRepository) *fileFacade {
	return &fileFacade{filesRepository: filesRepository}
}

func (ff *fileFacade) FindById(requesterId string, fileId string) (*entity.File, error) {
	return ff.filesRepository.FindById(requesterId, fileId)
}

func (ff *fileFacade) DeleteById(traceId string, requesterId string, fileId string) error {
	if err := ff.filesRepository.Delete(requesterId, fileId); err != nil {
		slog.Error("Could not delete file in database:", "traceId", traceId, "fileId", fileId, "error", err)
		return err
	}

	slog.Info("File removed successfully:", "traceId", traceId, "fileId", fileId)
	return nil
}

func (ff *fileFacade) FindAll(traceId string, requesterId string, page int, size int, filename string, secret bool) (*entity.FilePage, error) {
	if size == 0 || size > maxListSize {
		size = maxListSize
	}

	filesPage, err := ff.filesRepository.FindAll(requesterId, page, size, filename, secret)

	if err != nil {
		slog.Error("Could not list files", "traceId", traceId, "error", err)
		return nil, err
	}

	return filesPage, nil
}
