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
	IsEnabled     bool      `json:"isEnabled"`
	IsMfaEnabled  bool      `json:"isMfaEnabled"`
	IsMfaVerified bool      `json:"isMfaVerified"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

type ErrorResponse struct {
	Message string `json:"message,omitempty"`
	TraceId string `json:"traceId,omitempty"`
	Code    string `json:"code,omitempty"`
}

type PageRepresentation struct {
	Size          int                   `json:"size"`
	TotalElements int                   `json:"totalElements"`
	Page          int                   `json:"page"`
	Next          string                `json:"next,omitempty"`
	Content       []*UserRepresentation `json:"content"`
}
