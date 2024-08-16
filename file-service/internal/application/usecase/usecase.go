package usecase

import (
	"github.com/murilo-bracero/raspstore/file-service/internal/application/repository"
	"github.com/murilo-bracero/raspstore/file-service/internal/infra/config"
)

type UseCases struct {
	CreateFileUseCase   CreateFileUseCase
	UpdateFileUseCase   UpdateFileUseCase
	UploadUseCase       UploadFileUseCase
	DownloadFileUseCase DownloadFileUseCase
}

func InitUseCases(config *config.Config, repo repository.FilesRepository, txRepo repository.TxFilesRepository) *UseCases {
	return &UseCases{
		CreateFileUseCase:   NewCreateFileUseCase(config, repo),
		UpdateFileUseCase:   NewUpdateFileUseCase(txRepo),
		UploadUseCase:       NewUploadFileUseCase(config),
		DownloadFileUseCase: NewDownloadFileUseCase(config),
	}
}
