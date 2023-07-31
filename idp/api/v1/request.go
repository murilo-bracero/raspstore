package v1

type UpdateProfileRepresentation struct {
	Username string `json:"username,omitempty"`
}

type CreateUserRepresentation struct {
	Username string   `json:"username,omitempty"`
	Password string   `json:"password,omitempty"`
	Roles    []string `json:"roles,omitempty"`
}

type UpdateUserRepresentation struct {
	Enabled    bool     `json:"enabled,omitempty"`
	MfaEnabled bool     `json:"mfaEnabled,omitempty"`
	Username   string   `json:"username,omitempty"`
	Roles      []string `json:"roles,omitempty"`
}
