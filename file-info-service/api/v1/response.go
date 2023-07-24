package v1

import "github.com/murilo-bracero/raspstore/file-info-service/internal/model"

type FilePageResponse struct {
	Size          int                         `json:"size"`
	TotalElements int                         `json:"totalElements"`
	Page          int                         `json:"page"`
	Next          string                      `json:"next"`
	Content       []*model.FileMetadataLookup `json:"content"`
}

type ErrorResponse struct {
	Message string `json:"message,omitempty"`
	TraceId string `json:"traceId,omitempty"`
	Code    string `json:"code,omitempty"`
}
