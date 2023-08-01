package validators

import (
	v1 "github.com/murilo-bracero/raspstore/file-service/api/v1"
	"github.com/murilo-bracero/raspstore/file-service/internal"
)

func ValidateUpdateFileRequest(req *v1.UpdateFileRequest) error {
	if req.Editors == nil {
		return internal.ErrEditorsNil
	}

	if req.Viewers == nil {
		return internal.ErrViewersNil
	}

	if req.Filename == "" {
		return internal.ErrFilenameEmpty
	}

	return nil
}
