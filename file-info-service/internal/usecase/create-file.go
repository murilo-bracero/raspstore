package usecase

import (
	"github.com/murilo-bracero/raspstore/commons/pkg/logger"
	"github.com/murilo-bracero/raspstore/file-info-service/internal"
	"github.com/murilo-bracero/raspstore/file-info-service/internal/converter"
	"github.com/murilo-bracero/raspstore/file-info-service/internal/model"
	"github.com/murilo-bracero/raspstore/file-info-service/internal/repository"
)

type CreateFileUseCase interface {
	Execute(file *model.File) (error_ error)
}

type createFileUseCase struct {
	filesRepository repository.FilesRepository
}

func NewCreateFileUseCase(fr repository.FilesRepository) CreateFileUseCase {
	return &createFileUseCase{filesRepository: fr}
}

func (c *createFileUseCase) Execute(file *model.File) (error_ error) {
	usage, error_ := c.filesRepository.FindUsageByUserId(file.Owner)

	if error_ != nil {
		logger.Error("Could not find user usage: %s", error_.Error())
		return
	}

	if error_ != nil {
		logger.Error("Could not get user config: %s", error_.Error())
		return
	}

	available := int64(converter.ToIntBytes(internal.StorageLimit())) - usage

	if file.Size > available {
		logger.Info("Could not create file because available storage for user is insufficient: userId=%s, available=%d", file.Owner, available)
		return internal.ErrNotAvailableSpace
	}

	if error_ = c.filesRepository.Save(file); error_ != nil {
		logger.Error("Could not create file: %s", error_.Error())
		return
	}

	return
}
