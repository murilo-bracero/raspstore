package usecase

import (
	"context"
	"io"
	"log"
	"os"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/murilo-bracero/raspstore/file-info-service/proto/v1/file-info-service/pb"
	"raspstore.github.io/fs-service/internal"
	"raspstore.github.io/fs-service/internal/grpc/client"
)

type UploadFileUseCase interface {
	Execute(ctx context.Context, req *pb.CreateFileMetadataRequest, src io.Reader) (fileMetadata *pb.FileMetadata, error_ error)
}

type uploadFileUseCase struct {
	service client.FileInfoService
}

func NewUploadFileUseCase(service client.FileInfoService) UploadFileUseCase {
	return &uploadFileUseCase{service: service}
}

func (u *uploadFileUseCase) Execute(ctx context.Context, req *pb.CreateFileMetadataRequest, src io.Reader) (fileMetadata *pb.FileMetadata, error_ error) {
	traceId := ctx.Value(middleware.RequestIDKey).(string)

	fileMetadata, error_ = u.service.CreateFileMetadata(req)

	if error_ != nil {
		log.Printf("[ERROR] - [%s]: Could not create file in file-info-service due to error: %s", traceId, error_.Error())
		return
	}

	filerep, error_ := os.Create(internal.StoragePath() + "/" + fileMetadata.FileId)

	if error_ != nil {
		log.Printf("[ERROR] - [%s]: Could not create file in fs due to error: %s", traceId, error_.Error())
		return
	}

	defer filerep.Close()

	_, error_ = io.Copy(filerep, src)

	if error_ != nil {
		log.Printf("[ERROR] - [%s]: Could not read file buffer due to error: %s", traceId, error_.Error())
		return
	}

	return
}
