package dto

import "raspstore.github.io/file-manager/model"

type FileMetadataList struct {
	Size          int           `json:"size"`
	TotalElements int           `json:"totalElements"`
	Page          int           `json:"page"`
	Next          string        `json:"next"`
	Content       []*model.File `json:"content"`
}

type UpdateFileRequest struct {
	Path     string `json:"path,omitempty"`
	Filename string `json:"filename,omitempty"`
	Editors  string `json:"editors"`
	Viewers  string `json:"viewers"`
}

type ErrorResponse struct {
	Message string `json:"message,omitempty"`
	TraceId string `json:"traceId,omitempty"`
	Code    string `json:"code,omitempty"`
}
