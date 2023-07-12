package v1

type ErrorResponse struct {
	Message string `json:"message,omitempty"`
	TraceId string `json:"traceId,omitempty"`
	Code    string `json:"code,omitempty"`
}

type UserResponse struct {
	UserId      string `json:"userId,omitempty"`
	Username    string `json:"username,omitempty"`
	Email       string `json:"email,omitempty"`
	PhoneNumber string `json:"phoneNumber,omitempty"`
	CreatedAt   string `json:"createdAt,omitempty"`
	UpdatedAt   string `json:"updatedAt,omitempty"`
}

type AdminUserResponse struct {
	UserResponse
	Permissions []string `json:"permissions,omitempty"`
}

type UserResponseList struct {
	Page          int            `json:"page"`
	Size          int            `json:"size"`
	TotalElements int            `json:"totalElements"`
	Next          string         `json:"next"`
	Content       []UserResponse `json:"content"`
}

type UserConfigurationResponse struct {
	StorageLimit            string `json:"storageLimit,omitempty"`
	MinPasswordLength       int    `json:"minPasswordLength"`
	AllowPublicUserCreation bool   `json:"allowPublicUserCreation"`
	AllowLoginWithEmail     bool   `json:"allowLoginWithEmail"`
	EnforceMfa              bool   `json:"enforceMfa"`
}
