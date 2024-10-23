package facade

import (
	"errors"
	"log/slog"

	"github.com/murilo-bracero/raspstore/file-service/internal/application/parser"
	"github.com/murilo-bracero/raspstore/file-service/internal/domain/entity"
	"github.com/murilo-bracero/raspstore/file-service/internal/infra/config"
	"github.com/murilo-bracero/raspstore/file-service/internal/infra/repository"
)

const maxListSize = 50

var ErrNotAvailableSpace = errors.New("file is greather than the space available for your user")

type FileFacade interface {
	Save(file *entity.File) error
	FindById(requesterId string, fileId string) (*entity.File, error)
	DeleteById(traceId string, requesterId string, fileId string) error
	FindAll(traceId string, requesterId string, page int, size int, filename string, secret bool) (*entity.FilePage, error)
}

type fileFacade struct {
	config          *config.Config
	filesRepository repository.FilesRepository
}

func NewFileFacade(config *config.Config, filesRepository repository.FilesRepository) *fileFacade {
	return &fileFacade{config: config, filesRepository: filesRepository}
}

func (ff *fileFacade) Save(file *entity.File) error {
	usage, err := ff.filesRepository.FindUsageByUserId(file.Owner)

	if err != nil {
		slog.Error("Could not find user usage:", "error", err.Error())
		return err
	}

	available := int64(parser.ParseUsage(ff.config.Storage.Limit)) - usage

	if file.Size > available {
		slog.Info("Could not create file because available storage for user is insufficient:", "userId", file.Owner, "available", available)
		return ErrNotAvailableSpace
	}

	if err = ff.filesRepository.Save(file); err != nil {
		slog.Error("Could not create file:", "error", err.Error())
		return err
	}

	return err
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
