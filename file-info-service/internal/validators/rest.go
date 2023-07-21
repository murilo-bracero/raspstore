package validators

import (
	v1 "github.com/murilo-bracero/raspstore/file-info-service/api/v1"
	"github.com/murilo-bracero/raspstore/file-info-service/internal"
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

	if req.Folder == (v1.FolderRepresentation{}) {
		return internal.ErrFolderEmpty
	}

	if req.Folder.Name == "" {
		return internal.ErrFolderNameEmpty
	}

	return nil
}
