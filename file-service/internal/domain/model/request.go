package model

type UpdateFileRequest struct {
	Filename string `json:"filename,omitempty"`
	Secret   bool   `json:"secret"`
}
