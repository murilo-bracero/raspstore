package validators

import (
	"errors"

	"github.com/murilo-bracero/raspstore/file-service/internal/domain/model/request"
)

var (
	ErrFilenameEmpty = errors.New("field Filename must not be empty")
	ErrViewersNil    = errors.New("viewers field must be an array")
	ErrEditorsNil    = errors.New("editors field must be an array")
)

func ValidateUpdateFileRequest(req *request.UpdateFileRequest) error {
	if req.Editors == nil {
		return ErrEditorsNil
	}

	if req.Viewers == nil {
		return ErrViewersNil
	}

	if req.Filename == "" {
		return ErrFilenameEmpty
	}

	return nil
}
