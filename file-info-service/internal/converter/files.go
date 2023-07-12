package converter

import (
	"time"

	"github.com/google/uuid"
	"github.com/murilo-bracero/raspstore-protofiles/file-info-service/pb"
	"raspstore.github.io/file-manager/internal/model"
)

func ToFile(req *pb.CreateFileMetadataRequest) *model.File {
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

func ToFileMetadata(file *model.File) *pb.FileMetadata {
	return &pb.FileMetadata{
		FileId:   file.FileId,
		Filename: file.Filename,
		Path:     file.Path,
		OwnerId:  file.CreatedBy,
	}
}
