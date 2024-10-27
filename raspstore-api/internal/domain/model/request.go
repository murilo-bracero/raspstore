package model

type UpdateFileRequest struct {
	Filename string `json:"filename,omitempty"`
	Secret   bool   `json:"secret"`
}

type CreateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Name     string `json:"name,omitempty"`
}
