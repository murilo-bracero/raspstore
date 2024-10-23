package usecase

import (
	"log/slog"

	"github.com/murilo-bracero/raspstore/file-service/internal/domain/entity"
	"github.com/murilo-bracero/raspstore/file-service/internal/infra/repository"
)

type UpdateFileUseCase interface {
	Execute(file *entity.File, userId, traceId string) (fileMetadata *entity.File, err error)
}

type updateFileUseCase struct {
	repo repository.TxFilesRepository
}

func NewUpdateFileUseCase(repo repository.TxFilesRepository) *updateFileUseCase {
	return &updateFileUseCase{repo: repo}
}

func (c *updateFileUseCase) Execute(file *entity.File, userId, traceId string) (fileMetadata *entity.File, err error) {
	tx, err := c.repo.Begin()

	if err != nil {
		slog.Error("Could not initialize transaction", "traceId", traceId, "error", err)
		return nil, err
	}

	fileMetadata, err = c.repo.FindById(tx, userId, file.FileId)

	if err != nil {
		slog.Error("Could not find file", "traceId", traceId, "fileId", file.FileId, "error", err)
		return
	}

	fileMetadata.Secret = file.Secret
	fileMetadata.Filename = file.Filename

	if fileMetadata.Secret {
		err := c.repo.DeleteFilePermissionByFileId(tx, fileMetadata.FileId)
		if err != nil {
			slog.Error("Could not remove permissions to set file secret", "traceId", traceId, "fileId", fileMetadata.FileId, "error", err)
			return nil, err
		}
	}

	if err = c.repo.Update(tx, userId, fileMetadata); err != nil {
		slog.Error("Could not update file", "traceId", traceId, "fileId", file.FileId, "error", err)
		return nil, err
	}

	if err := c.repo.Commit(tx); err != nil {
		slog.Error("Could commit update file transaction", "traceId", traceId, "fileId", file.FileId, "error", err)
	}

	slog.Info("File updated successfully", "traceId", traceId, "fileId", file.FileId)
	return
}
