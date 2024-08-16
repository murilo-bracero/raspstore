package repository

import (
	"database/sql"
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
	DeleteFilePermissionByFileId(fileId string) error
}

type TxFilesRepository interface {
	Begin() (*sql.Tx, error)
	Commit(tx *sql.Tx) error
	FindById(tx *sql.Tx, userId string, fileId string) (*entity.File, error)
	Update(tx *sql.Tx, userId string, file *entity.File) error
	DeleteFilePermissionByFileId(tx *sql.Tx, fileId string) error
}
