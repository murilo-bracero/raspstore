package usecase

import (
	"github.com/murilo-bracero/raspstore/file-service/internal/application/repository"
	"github.com/murilo-bracero/raspstore/file-service/internal/domain/entity"
)

type getFileUseCase struct {
	repo repository.FilesRepository
}

type GetFileUseCase interface {
	Execute(userId string, fileId string) (file *entity.File, error_ error)
}

func NewGetFileUseCase(repo repository.FilesRepository) GetFileUseCase {
	return &getFileUseCase{repo: repo}
}

func (c *getFileUseCase) Execute(userId string, fileId string) (file *entity.File, error_ error) {
	file, error_ = c.repo.FindById(userId, fileId)

	return
}