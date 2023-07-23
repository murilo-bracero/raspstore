package validators

import (
	"github.com/murilo-bracero/raspstore/file-info-service/internal"
	"github.com/murilo-bracero/raspstore/file-info-service/proto/v1/file-info-service/pb"
)

func ValidateCreateFileMetadataRequest(req *pb.CreateFileMetadataRequest) error {
	if req.Filename == "" {
		return internal.ErrFilenameEmpty
	}

	if req.OwnerId == "" {
		return internal.ErrOwnerIdEmpty
	}

	if req.Size <= 0 {
		return internal.ErrInvalidSize
	}

	return nil
}

func ValidateFindFileMetadataByIdRequest(req *pb.FindFileMetadataByIdRequest) error {
	if req.FileId == "" {
		return internal.ErrEmptyFileId
	}

	return nil
}
