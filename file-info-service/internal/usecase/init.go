package usecase

import (
	"github.com/murilo-bracero/raspstore/file-info-service/internal/repository"
)

type useCases struct {
	CreateFileUseCase CreateFileUseCase
	DeleteFileUseCase DeleteFileUseCase
	ListFilesUseCase  ListFilesUseCase
	UpdateFileUseCase UpdateFileUseCase
	GetFileUseCase    GetFileUseCase
	UploadUseCase     UploadFileUseCase
}

func InitUseCases(repo repository.FilesRepository) *useCases {
	return &useCases{
		CreateFileUseCase: NewCreateFileUseCase(repo),
		DeleteFileUseCase: NewDeleteFileUseCase(repo),
		ListFilesUseCase:  NewListFilesUseCase(repo),
		UpdateFileUseCase: NewUpdateFileUseCase(repo),
		GetFileUseCase:    NewGetFileUseCase(repo),
		UploadUseCase:     NewUploadFileUseCase(),
	}
}
