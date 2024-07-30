package usecase

import (
	"errors"
	"log/slog"

	"github.com/murilo-bracero/raspstore/file-service/internal/application/parser"
	"github.com/murilo-bracero/raspstore/file-service/internal/application/repository"
	"github.com/murilo-bracero/raspstore/file-service/internal/domain/entity"
	"github.com/murilo-bracero/raspstore/file-service/internal/infra/config"
)

var (
	ErrNotAvailableSpace = errors.New("file is greather than the space available for your user")
)

type CreateFileUseCase interface {
	Execute(file *entity.File) (err error)
}

type createFileUseCase struct {
	config          *config.Config
	filesRepository repository.FilesRepository
}

func NewCreateFileUseCase(config *config.Config, fr repository.FilesRepository) *createFileUseCase {
	return &createFileUseCase{filesRepository: fr, config: config}
}

func (c *createFileUseCase) Execute(file *entity.File) (err error) {
	usage, err := c.filesRepository.FindUsageByUserId(file.Owner)

	if err != nil {
		slog.Error("Could not find user usage:", "error", err.Error())
		return
	}

	available := int64(parser.ParseUsage(c.config.Storage.Limit)) - usage

	if file.Size > available {
		slog.Info("Could not create file because available storage for user is insufficient:", "userId", file.Owner, "available", available)
		return ErrNotAvailableSpace
	}

	if err = c.filesRepository.Save(file); err != nil {
		slog.Error("Could not create file:", "error", err.Error())
		return
	}

	return
}
