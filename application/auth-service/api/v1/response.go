package v1

import (
	"encoding/json"
	"net/http"
)

type LoginResponse struct {
	AccessToken           string `json:"accessToken,omitempty"`
	AccessTokenExpiresIn  int64  `json:"accessTokenExpiresIn,omitempty"`
	RefreshToken          string `json:"refreshToken,omitempty"`
	RefreshTokenExpiresIn int64  `json:"refreshTokenExpiresIn,omitempty"`
}

func Send(w http.ResponseWriter, obj interface{}) {
	w.Header().Set("Content-Type", "application/json")
	jsonResponse, err := json.Marshal(obj)
	if err != nil {
		return
	}
	w.Write(jsonResponse)
}
