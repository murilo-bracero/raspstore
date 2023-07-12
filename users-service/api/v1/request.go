package v1

type CreateUserRequest struct {
	Username    string `json:"username,omitempty"`
	Email       string `json:"email,omitempty"`
	Password    string `json:"password,omitempty"`
	PhoneNumber string `json:"phoneNumber,omitempty"`
}

type UpdateUserRequest struct {
	Username    string `json:"username,omitempty"`
	Email       string `json:"email,omitempty"`
	PhoneNumber string `json:"phoneNumber,omitempty"`
}

type PatchUserConfigRequest struct {
	StorageLimit            *string `json:"storageLimit,omitempty"`
	MinPasswordLength       *int    `json:"minPasswordLength,omitempty"`
	AllowPublicUserCreation *bool   `json:"allowPublicUserCreation,omitempty"`
	AllowLoginWithEmail     *bool   `json:"allowLoginWithEmail,omitempty"`
	EnforceMfa              *bool   `json:"enforceMfa,omitempty"`
}

type AdminCreateUserRequest struct {
	Permissions []string `json:"permissions,omitempty"`
	CreateUserRequest
}
