package v1

type LoginRequest struct {
	MfaToken string `json:"mfaToken"`
}
