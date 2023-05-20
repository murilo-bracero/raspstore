package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/murilo-bracero/raspstore-protofiles/file-info-service/pb"
	"raspstore.github.io/file-manager/internal"
	"raspstore.github.io/file-manager/model"
	"raspstore.github.io/file-manager/repository"
)

type fileInfoService struct {
	fileRepository repository.FilesRepository
	pb.UnimplementedFileInfoServiceServer
}

func NewFileInfoService(fileRepository repository.FilesRepository) pb.FileInfoServiceServer {
	return &fileInfoService{fileRepository: fileRepository}
}

func (f *fileInfoService) CreateFileMetadata(ctx context.Context, req *pb.CreateFileMetadataRequest) (*pb.FileMetadata, error) {

	if req.Filename == "" {
		return nil, internal.ErrFilenameEmpty
	}

	if req.OwnerId == "" {
		return nil, internal.ErrOwnerIdEmpty
	}

	if req.Path == "" {
		return nil, internal.ErrPathEmpty
	}

	if req.Size <= 0 {
		return nil, internal.ErrInvalidSize
	}

	file := &model.File{
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

	if err := f.fileRepository.Save(file); err != nil {
		return nil, err
	}

	return &pb.FileMetadata{
		FileId:   file.FileId,
		Filename: file.Filename,
		Path:     file.Path,
		OwnerId:  file.CreatedBy,
	}, nil
}

func (f *fileInfoService) FindFileMetadataById(ctx context.Context, req *pb.FindFileMetadataByIdRequest) (*pb.FileMetadata, error) {

	if req.FileId == "" {
		return nil, internal.ErrEmptyFileId
	}

	file, err := f.fileRepository.FindById(req.RequesterUserId, req.FileId)

	if err != nil {
		return nil, err
	}

	return &pb.FileMetadata{
		FileId:   file.FileId,
		Filename: file.Filename,
		Path:     file.Path,
		OwnerId:  file.CreatedBy,
	}, nil
}
