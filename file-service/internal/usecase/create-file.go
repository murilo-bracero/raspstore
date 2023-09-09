package usecase

import (
	"log/slog"

	"github.com/murilo-bracero/raspstore/file-service/internal"
	"github.com/murilo-bracero/raspstore/file-service/internal/converter"
	"github.com/murilo-bracero/raspstore/file-service/internal/infra"
	"github.com/murilo-bracero/raspstore/file-service/internal/model"
	"github.com/murilo-bracero/raspstore/file-service/internal/repository"
)

type CreateFileUseCase interface {
	Execute(file *model.File) (error_ error)
}

type createFileUseCase struct {
	config          *infra.Config
	filesRepository repository.FilesRepository
}

func NewCreateFileUseCase(config *infra.Config, fr repository.FilesRepository) CreateFileUseCase {
	return &createFileUseCase{filesRepository: fr, config: config}
}

func (c *createFileUseCase) Execute(file *model.File) (error_ error) {
	usage, error_ := c.filesRepository.FindUsageByUserId(file.Owner)

	if error_ != nil {
		slog.Error("Could not find user usage:", "error", error_.Error())
		return
	}

	if error_ != nil {
		slog.Error("Could not get user config:", "error", error_)
		return
	}

	available := int64(converter.ToIntBytes(c.config.Server.Storage.Limit)) - usage

	if file.Size > available {
		slog.Info("Could not create file because available storage for user is insufficient:", "userId", file.Owner, "available", available)
		return internal.ErrNotAvailableSpace
	}

	if error_ = c.filesRepository.Save(file); error_ != nil {
		slog.Error("Could not create file:", "error", error_.Error())
		return
	}

	return
}
