package validators

import (
	"github.com/murilo-bracero/raspstore-protofiles/file-info-service/pb"
	"raspstore.github.io/file-manager/internal"
)

func ValidateCreateFileMetadataRequest(req *pb.CreateFileMetadataRequest) error {
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

func ValidateFindFileMetadataByIdRequest(req *pb.FindFileMetadataByIdRequest) error {
	if req.FileId == "" {
		return internal.ErrEmptyFileId
	}

	return nil
}
