package v1

type CreateUserRequest struct {
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

type UpdateUserRequest struct {
	Username  string `json:"username,omitempty"`
	Email     string `json:"email,omitempty"`
	IsEnabled *bool  `json:"isEnabled,omitempty"`
	Password  string `json:"password,omitempty"`
}
