package model

type UpdateFileRequest struct {
	Filename string   `json:"filename,omitempty"`
	Secret   bool     `json:"secret"`
	Editors  []string `json:"editors"`
	Viewers  []string `json:"viewers"`
}
