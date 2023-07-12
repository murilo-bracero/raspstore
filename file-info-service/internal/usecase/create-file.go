package usecase

import (
	"log"

	"github.com/murilo-bracero/raspstore/file-info-service/internal"
	"github.com/murilo-bracero/raspstore/file-info-service/internal/converter"
	"github.com/murilo-bracero/raspstore/file-info-service/internal/grpc/client"
	"github.com/murilo-bracero/raspstore/file-info-service/internal/model"
	"github.com/murilo-bracero/raspstore/file-info-service/internal/repository"
)

type CreateFileUseCase interface {
	Execute(file *model.File) (error_ error)
}

type createFileUseCase struct {
	filesRepository   repository.FilesRepository
	userConfigService client.UserConfigGrpcService
}

func NewCreateFileUseCase(fr repository.FilesRepository, ucs client.UserConfigGrpcService) CreateFileUseCase {
	return &createFileUseCase{filesRepository: fr, userConfigService: ucs}
}

func (c *createFileUseCase) Execute(file *model.File) (error_ error) {
	usage, error_ := c.filesRepository.FindUsageByUserId(file.Owner)

	if error_ != nil {
		log.Printf("[ERROR] Could not find user usage: %s", error_.Error())
		return
	}

	userConfig, error_ := c.userConfigService.GetUserConfiguration()

	if error_ != nil {
		log.Printf("[ERROR] Could not get user config: %s", error_.Error())
		return
	}

	available := int64(converter.ToIntBytes(userConfig.StorageLimit)) - usage

	if file.Size > available {
		log.Printf("[INFO] Could not create file because available storage for user is insufficient: userId=%s, available=%d", file.Owner, available)
		return internal.ErrNotAvailableSpace
	}

	if error_ = c.filesRepository.Save(file); error_ != nil {
		log.Printf("[ERROR] Could not create file: %s", error_.Error())
		return
	}

	return
}
