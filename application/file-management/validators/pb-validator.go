package validators

import (
	"errors"

	"raspstore.github.io/file-manager/pb"
)

var (
	ErrReceiveStreaming = errors.New("an error occured while receiving the stream, try again later")
	ErrUploadFile       = errors.New("an error occured while writing uploaded file in server. try again later")
	ErrCreatedByEmpty   = errors.New("filed CreatedBy must not be empty")
	ErrFilenameEmpty    = errors.New("field Filename must not be empty")
	ErrWrongID          = errors.New("provided id is invalid")
)

func ValidateDownload(req *pb.GetFileRequest) error {
	if len(req.Id) != 24 {
		return ErrWrongID
	}

	return nil
}

func ValidateUpload(req *pb.CreateFileRequest) error {

	if req.GetFiledata().CreatedBy == "" {
		return ErrCreatedByEmpty
	}

	if req.GetFiledata().Filename == "" {
		return ErrFilenameEmpty
	}

	return nil
}

func ValidateUpdate(req *pb.UpdateFileRequest) error {
	if len(req.GetFiledata().Id) != 24 {
		return ErrWrongID
	}

	if req.GetFiledata().UpdatedBy == "" {
		return ErrCreatedByEmpty
	}

	return nil
}
