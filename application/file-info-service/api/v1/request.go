package v1

type UpdateFileRequest struct {
	Path     string   `json:"path,omitempty"`
	Filename string   `json:"filename,omitempty"`
	Editors  []string `json:"editors"`
	Viewers  []string `json:"viewers"`
}
