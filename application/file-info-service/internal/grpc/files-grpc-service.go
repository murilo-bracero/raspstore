package grpc

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/murilo-bracero/raspstore-protofiles/file-info-service/pb"
	"raspstore.github.io/file-manager/internal"
	"raspstore.github.io/file-manager/internal/model"
	"raspstore.github.io/file-manager/internal/repository"
)

type fileInfoService struct {
	fileRepository repository.FilesRepository
	pb.UnimplementedFileInfoServiceServer
}

func NewFileInfoService(fileRepository repository.FilesRepository) pb.FileInfoServiceServer {
	return &fileInfoService{fileRepository: fileRepository}
}

func (f *fileInfoService) CreateFileMetadata(ctx context.Context, req *pb.CreateFileMetadataRequest) (*pb.FileMetadata, error) {

	if err := validateCreateFileMetadataRequest(req); err != nil {
		return nil, err
	}

	file := buildFile(req)

	if err := f.fileRepository.Save(file); err != nil {
		return nil, err
	}

	return buildFileMetadata(file), nil
}

func (f *fileInfoService) FindFileMetadataById(ctx context.Context, req *pb.FindFileMetadataByIdRequest) (*pb.FileMetadata, error) {

	if err := validateFindFileMetadataByIdRequest(req); err != nil {
		return nil, err
	}

	file, err := f.fileRepository.FindById(req.RequesterUserId, req.FileId)

	if err != nil {
		return nil, err
	}

	return buildFileMetadata(file), nil
}

func validateCreateFileMetadataRequest(req *pb.CreateFileMetadataRequest) error {
	if req.Filename == "" {
		return internal.ErrFilenameEmpty
	}

	if req.OwnerId == "" {
		return internal.ErrOwnerIdEmpty
	}

	if req.Path == "" {
		return internal.ErrPathEmpty
	}

	if req.Size <= 0 {
		return internal.ErrInvalidSize
	}

	return nil
}

func validateFindFileMetadataByIdRequest(req *pb.FindFileMetadataByIdRequest) error {
	if req.FileId == "" {
		return internal.ErrEmptyFileId
	}

	return nil
}

func buildFile(req *pb.CreateFileMetadataRequest) *model.File {
	return &model.File{
		FileId:    uuid.NewString(),
		Filename:  req.Filename,
		Path:      req.Path,
		Size:      req.Size,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Viewers:   []string{},
		Editors:   []string{},
		CreatedBy: req.OwnerId,
		UpdatedBy: req.OwnerId,
		Owner:     req.OwnerId,
	}
}

func buildFileMetadata(file *model.File) *pb.FileMetadata {
	return &pb.FileMetadata{
		FileId:   file.FileId,
		Filename: file.Filename,
		Path:     file.Path,
		OwnerId:  file.CreatedBy,
	}
}
