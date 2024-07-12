package usecase

import (
	"github.com/murilo-bracero/raspstore/file-service/internal/application/repository"
	"github.com/murilo-bracero/raspstore/file-service/internal/infra/config"
)

type UseCases struct {
	CreateFileUseCase   CreateFileUseCase
	DeleteFileUseCase   DeleteFileUseCase
	ListFilesUseCase    ListFilesUseCase
	UpdateFileUseCase   UpdateFileUseCase
	GetFileUseCase      GetFileUseCase
	UploadUseCase       UploadFileUseCase
	DownloadFileUseCase DownloadFileUseCase
}

func InitUseCases(config *config.Config, repo repository.FilesRepository) *UseCases {
	return &UseCases{
		CreateFileUseCase:   NewCreateFileUseCase(config, repo),
		DeleteFileUseCase:   NewDeleteFileUseCase(repo),
		ListFilesUseCase:    NewListFilesUseCase(repo),
		UpdateFileUseCase:   NewUpdateFileUseCase(repo),
		GetFileUseCase:      NewGetFileUseCase(repo),
		UploadUseCase:       NewUploadFileUseCase(config),
		DownloadFileUseCase: NewDownloadFileUseCase(config),
	}
}
