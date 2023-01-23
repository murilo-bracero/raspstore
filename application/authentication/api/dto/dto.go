package dto

type LoginRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	TotpToken string `json:"totpToken"`
}

type LoginResponse struct {
	Token string `json:"token,omitempty"`
}
