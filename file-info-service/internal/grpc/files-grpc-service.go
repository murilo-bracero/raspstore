package grpc

import (
	"context"

	"github.com/murilo-bracero/raspstore-protofiles/file-info-service/pb"
	"raspstore.github.io/file-manager/internal/converter"
	"raspstore.github.io/file-manager/internal/repository"
	"raspstore.github.io/file-manager/internal/validators"
)

type fileInfoService struct {
	fileRepository repository.FilesRepository
	pb.UnimplementedFileInfoServiceServer
}

func NewFileInfoService(fileRepository repository.FilesRepository) pb.FileInfoServiceServer {
	return &fileInfoService{fileRepository: fileRepository}
}

func (f *fileInfoService) CreateFileMetadata(ctx context.Context, req *pb.CreateFileMetadataRequest) (*pb.FileMetadata, error) {

	if err := validators.ValidateCreateFileMetadataRequest(req); err != nil {
		return nil, err
	}

	file := converter.ToFile(req)

	if err := f.fileRepository.Save(file); err != nil {
		return nil, err
	}

	return converter.ToFileMetadata(file), nil
}

func (f *fileInfoService) FindFileMetadataById(ctx context.Context, req *pb.FindFileMetadataByIdRequest) (*pb.FileMetadata, error) {

	if err := validators.ValidateFindFileMetadataByIdRequest(req); err != nil {
		return nil, err
	}

	file, err := f.fileRepository.FindById(req.RequesterUserId, req.FileId)

	if err != nil {
		return nil, err
	}

	return converter.ToFileMetadata(file), nil
}
