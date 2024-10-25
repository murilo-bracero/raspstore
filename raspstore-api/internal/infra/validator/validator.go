package validator

import (
	"errors"

	"github.com/murilo-bracero/raspstore/file-service/internal/domain/model"
)

var ErrFilenameEmpty = errors.New("field Filename must not be empty")

func ValidateUpdateFileRequest(req *model.UpdateFileRequest) error {
	if req.Filename == "" {
		return ErrFilenameEmpty
	}

	return nil
}
