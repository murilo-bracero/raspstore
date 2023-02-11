package validators

import (
	"strings"

	"raspstore.github.io/file-manager/internal"
)

func ValidateBody(received string, desired string) bool {
	clean := strings.Split(received, ";")

	return clean[0] == desired
}

func ValidateId(id string) error {
	if len(id) != 24 {
		return internal.ErrWrongID
	}

	return nil
}
