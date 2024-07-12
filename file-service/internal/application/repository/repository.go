package repository

import (
	"errors"

	"github.com/murilo-bracero/raspstore/file-service/internal/domain/entity"
)

var ErrFileDoesNotExists = errors.New("file with provided ID does not exists")

type FilesRepository interface {
	Save(file *entity.File) error
	FindById(userId string, fileId string) (*entity.File, error)
	FindUsageByUserId(userId string) (usage int64, err error)
	Delete(userId string, fileId string) error
	Update(userId string, file *entity.File) error
	FindAll(userId string, page int, size int, filename string, secret bool) (filesPage *entity.FilePage, err error)
}
