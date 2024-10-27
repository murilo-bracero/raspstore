package entity

import "github.com/google/uuid"

type User struct {
	Id           uuid.UUID
	Username     string
	PasswordHash string
	RefreshToken string
	Name         string
}

func NewUser(username, passwordHash, name, refreshToken string) *User {
	return &User{
		Id:           uuid.New(),
		Username:     username,
		PasswordHash: passwordHash,
		Name:         name,
		RefreshToken: refreshToken,
	}
}
