package model

import "time"

type TokenCredentials struct {
	AccessToken  string
	RefreshToken string
	ExpirestAt   time.Time
}
