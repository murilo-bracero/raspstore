package usecase

import (
	"github.com/murilo-bracero/raspstore/file-service/internal/infra/config"
	"github.com/murilo-bracero/raspstore/file-service/internal/infra/repository"
)

type UseCases struct {
	CreateFileUseCase   CreateFileUseCase
	UpdateFileUseCase   UpdateFileUseCase
	UploadUseCase       UploadFileUseCase
	DownloadFileUseCase DownloadFileUseCase
	LoginPAMUseCase     LoginPAMUseCase
}

func InitUseCases(config *config.Config, repo repository.FilesRepository, txRepo repository.TxFilesRepository) *UseCases {
	return &UseCases{
		CreateFileUseCase:   NewCreateFileUseCase(config, repo),
		UpdateFileUseCase:   NewUpdateFileUseCase(txRepo),
		UploadUseCase:       NewUploadFileUseCase(config),
		DownloadFileUseCase: NewDownloadFileUseCase(config),
		LoginPAMUseCase:     NewLoginPAMUseCase(config),
	}
}
