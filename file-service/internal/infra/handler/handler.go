package handler

import (
	"github.com/murilo-bracero/raspstore/file-service/internal/application/facade"
	"github.com/murilo-bracero/raspstore/file-service/internal/application/usecase"
	"github.com/murilo-bracero/raspstore/file-service/internal/infra/config"
)

type Handler struct {
	downloadUseCase   usecase.DownloadFileUseCase
	uploadUseCase     usecase.UploadFileUseCase
	createFileUseCase usecase.CreateFileUseCase
	updateFileUseCase usecase.UpdateFileUseCase
	loginPAMUseCase   usecase.LoginPAMUseCase
	fileFacade        facade.FileFacade
	config            *config.Config
}

func New(
	downloadUseCase usecase.DownloadFileUseCase,
	uploadUseCase usecase.UploadFileUseCase,
	createFileUseCase usecase.CreateFileUseCase,
	updateFileUseCase usecase.UpdateFileUseCase,
	loginPAMUseCase usecase.LoginPAMUseCase,
	fileFacade facade.FileFacade,
	config *config.Config,
) *Handler {
	return &Handler{downloadUseCase, uploadUseCase, createFileUseCase, updateFileUseCase, loginPAMUseCase, fileFacade, config}
}
