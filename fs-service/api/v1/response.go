package v1

type UploadSuccessResponse struct {
	FileId   string `json:"fileId,omitempty"`
	Filename string `json:"filename,omitempty"`
	Path     string `json:"path,omitempty"`
	OwnerId  string `json:"ownerId,omitempty"`
}

type ErrorResponse struct {
	Message string `json:"message,omitempty"`
	TraceId string `json:"traceId,omitempty"`
	Code    string `json:"code,omitempty"`
}
