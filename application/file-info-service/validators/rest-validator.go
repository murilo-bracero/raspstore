package validators

import "strings"

func ValidateBody(received string, desired string) bool {
	clean := strings.Split(received, ";")

	return clean[0] == desired
}

func ValidateId(id string) error {
	if len(id) != 24 {
		return ErrWrongID
	}

	return nil
}
