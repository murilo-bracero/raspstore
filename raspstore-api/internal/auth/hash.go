package auth

import "golang.org/x/crypto/bcrypt"

func HashPassword(raw string) (string, error) {
	digest, err := bcrypt.GenerateFromPassword([]byte(raw), bcrypt.DefaultCost)

	return string(digest), err
}
