package model

import (
	"time"
)

type FilePageResponse struct {
	Size          int            `json:"size"`
	TotalElements int            `json:"totalElements"`
	Page          int            `json:"page"`
	Next          string         `json:"next"`
	Content       []*FileContent `json:"content"`
}

type FileContent struct {
	FileId    string     `json:"fileId,omitempty"`
	Filename  string     `json:"filename,omitempty"`
	Size      int64      `json:"size,omitempty"`
	Owner     string     `json:"owner,omitempty"`
	CreatedAt time.Time  `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
	CreatedBy string     `json:"createdBy,omitempty"`
	UpdatedBy *string    `json:"updatedBy,omitempty"`
}

type ErrorResponse struct {
	Message string `json:"message,omitempty"`
	TraceId string `json:"traceId,omitempty"`
	Code    string `json:"code,omitempty"`
}

type UploadSuccessResponse struct {
	FileId   string `json:"fileId,omitempty"`
	Filename string `json:"filename,omitempty"`
	OwnerId  string `json:"ownerId,omitempty"`
}
