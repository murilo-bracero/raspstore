package converter

import (
	"time"

	"github.com/google/uuid"
	"github.com/murilo-bracero/raspstore/file-info-service/internal/model"
	"github.com/murilo-bracero/raspstore/file-info-service/proto/v1/file-info-service/pb"
)

func ToFile(req *pb.CreateFileMetadataRequest) *model.File {
	return &model.File{
		FileId:    uuid.NewString(),
		Filename:  req.Filename,
		Size:      req.Size,
		Secret:    req.Secret,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Viewers:   []string{},
		Editors:   []string{},
		CreatedBy: req.OwnerId,
		UpdatedBy: req.OwnerId,
		Owner:     req.OwnerId,
	}
}

func ToFileMetadata(file *model.File) *pb.FileMetadata {
	return &pb.FileMetadata{
		FileId:   file.FileId,
		Filename: file.Filename,
		OwnerId:  file.CreatedBy,
	}
}
