package v1

import "time"

type LoginResponse struct {
	AccessToken           string `json:"accessToken,omitempty"`
	AccessTokenExpiresIn  int64  `json:"accessTokenExpiresIn,omitempty"`
	RefreshToken          string `json:"refreshToken,omitempty"`
	RefreshTokenExpiresIn int64  `json:"refreshTokenExpiresIn,omitempty"`
}

type UserRepresentation struct {
	UserID        string    `json:"userId"`
	Username      string    `json:"username"`
	Email         string    `json:"email"`
	PhoneNumber   string    `json:"phoneNumber"`
	IsMfaEnabled  bool      `json:"isMfaEnabled"`
	IsMfaVerified bool      `json:"isMfaVerified"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}
