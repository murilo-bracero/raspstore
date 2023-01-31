package dto

type LoginRequest struct {
	MfaToken string `json:"mfaToken"`
}

type LoginResponse struct {
	AccessToken           string `json:"accessToken,omitempty"`
	AccessTokenExpiresIn  int64  `json:"accessTokenExpiresIn,omitempty"`
	RefreshToken          string `json:"refreshToken,omitempty"`
	RefreshTokenExpiresIn int64  `json:"refreshTokenExpiresIn,omitempty"`
}
