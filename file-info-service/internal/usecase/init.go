package usecase

import (
	"github.com/murilo-bracero/raspstore/file-info-service/internal/grpc/client"
	"github.com/murilo-bracero/raspstore/file-info-service/internal/repository"
)

type useCases struct {
	CreateFileUseCase CreateFileUseCase
	DeleteFileUseCase DeleteFileUseCase
	ListFilesUseCase  ListFilesUseCase
	UpdateFileUseCase UpdateFileUseCase
	GetFileUseCase    GetFileUseCase
}

func InitUseCases(repo repository.FilesRepository, ucs client.UserConfigGrpcService) *useCases {
	return &useCases{
		CreateFileUseCase: NewCreateFileUseCase(repo, ucs),
		DeleteFileUseCase: NewDeleteFileUseCase(repo),
		ListFilesUseCase:  NewListFilesUseCase(repo),
		UpdateFileUseCase: NewUpdateFileUseCase(repo),
		GetFileUseCase:    NewGetFileUseCase(repo),
	}
}