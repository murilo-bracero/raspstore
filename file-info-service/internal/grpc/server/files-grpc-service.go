package server

import (
	"context"

	"github.com/murilo-bracero/raspstore/file-info-service/internal/converter"
	"github.com/murilo-bracero/raspstore/file-info-service/internal/usecase"
	"github.com/murilo-bracero/raspstore/file-info-service/internal/validators"
	"github.com/murilo-bracero/raspstore/file-info-service/proto/v1/file-info-service/pb"
)

type fileInfoService struct {
	getFileUseCase    usecase.GetFileUseCase
	createFileUseCase usecase.CreateFileUseCase
	pb.UnimplementedFileInfoServiceServer
}

func NewFileInfoService(getFileUseCase usecase.GetFileUseCase, cfuc usecase.CreateFileUseCase) pb.FileInfoServiceServer {
	return &fileInfoService{getFileUseCase: getFileUseCase, createFileUseCase: cfuc}
}

func (f *fileInfoService) CreateFileMetadata(ctx context.Context, req *pb.CreateFileMetadataRequest) (*pb.FileMetadata, error) {

	if err := validators.ValidateCreateFileMetadataRequest(req); err != nil {
		return nil, err
	}

	file := converter.ToFile(req)

	if err := f.createFileUseCase.Execute(file); err != nil {
		return nil, err
	}

	return converter.ToFileMetadata(file), nil
}

func (f *fileInfoService) FindFileMetadataById(ctx context.Context, req *pb.FindFileMetadataByIdRequest) (*pb.FileMetadata, error) {

	if err := validators.ValidateFindFileMetadataByIdRequest(req); err != nil {
		return nil, err
	}

	file, err := f.getFileUseCase.Execute(req.RequesterUserId, req.FileId)

	if err != nil {
		return nil, err
	}

	return converter.ToFileMetadata(file), nil
}
