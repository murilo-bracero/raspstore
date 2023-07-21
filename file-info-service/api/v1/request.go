package v1

type UpdateFileRequest struct {
	Folder   FolderRepresentation `json:"folder,omitempty"`
	Filename string               `json:"filename,omitempty"`
	Editors  []string             `json:"editors"`
	Viewers  []string             `json:"viewers"`
}

type FolderRepresentation struct {
	Name   string `json:"name,omitempty"`
	Secret bool   `json:"secret,omitempty"`
}
