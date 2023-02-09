package service

import (
	"context"

	"github.com/murilo-bracero/raspstore-protofiles/file-info-service/pb"
	"raspstore.github.io/file-manager/repository"
)

type fileInfoService struct {
	fileRepository repository.FilesRepository
	pb.UnimplementedFileInfoServiceServer
}

func NewFileManagerService(fileRepository repository.FilesRepository) pb.FileInfoServiceServer {
	return &fileInfoService{fileRepository: fileRepository}
}

func (f *fileInfoService) CreateFileMetadata(ctx context.Context, req *pb.CreateFileMetadataRequest) (*pb.FileMetadata, error) {

	return &pb.FileMetadata{}, nil
}

func (f *fileInfoService) FindFileMetadataById(ctx context.Context, req *pb.FindFileMetadataByIdRequest) (*pb.FileMetadata, error) {

	return &pb.FileMetadata{}, nil
}
