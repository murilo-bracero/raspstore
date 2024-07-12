package usecase

import (
	"github.com/murilo-bracero/raspstore/file-service/internal/application/repository"
	"github.com/murilo-bracero/raspstore/file-service/internal/domain/entity"
)

type getFileUseCase struct {
	repo repository.FilesRepository
}

type GetFileUseCase interface {
	Execute(userId string, fileId string) (file *entity.File, err error)
}

func NewGetFileUseCase(repo repository.FilesRepository) *getFileUseCase {
	return &getFileUseCase{repo: repo}
}

func (c *getFileUseCase) Execute(userId string, fileId string) (file *entity.File, err error) {
	file, err = c.repo.FindById(userId, fileId)

	return
}
