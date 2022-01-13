package dto

type FileMetadata struct {
	Uri       string `json:"uri,omitempty"`
	UpdatedAt int    `json:"updatedAt,omitempty"`
	CreatedBy string `json:"createdBy,omitempty"`
	UpdatedBy string `json:"updatedBy,omitempty"`
}

type CreateFileRequestData struct {
	Filename  string `json:"filename,omitempty"`
	CreatedBy string `json:"createdBy,omitempty"`
}

type UpdateFileRequestData struct {
	Id        string `json:"id,omitempty"`
	UpdatedBy string `json:"updatedBy,omitempty"`
}

type CreateFileRequest struct {
	Filedata CreateFileRequestData `json:"id,omitempty"`
	Chunk    []byte                `json:"chunk,omitempty"`
}

type UpdateFileRequest struct {
	Filedata UpdateFileRequestData `json:"filedata,omitempty"`
	Chunk    []byte                `json:"chunk,omitempty"`
}

type ErrorResponse struct {
	Message string `json:"message,omitempty"`
	Reason  string `json:"reason,omitempty"`
	Code    string `json:"code,omitempty"`
}
