package usecase

import (
	"log"

	"raspstore.github.io/file-manager/internal/model"
	"raspstore.github.io/file-manager/internal/repository"
)

type CreateFileUseCase interface {
	Execute(file *model.File) (result *model.File, error_ error)
}

type createFileUseCase struct {
	filesRepository repository.FilesRepository
}

func NewCreateFileUseCase(filesRepository repository.FilesRepository) CreateFileUseCase {
	return &createFileUseCase{filesRepository: filesRepository}
}

func (c *createFileUseCase) Execute(file *model.File) (result *model.File, error_ error) {
	usage, error_ := c.filesRepository.FindUserUsageById(file.Owner)

	if error_ != nil {
		log.Printf("[ERROR] Could not create file: %s", error_.Error())
		return
	}

	if error_ = c.filesRepository.Save(file); error_ != nil {
		return
	}

	return
}
