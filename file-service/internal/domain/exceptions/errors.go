package exceptions

import "errors"

var (
	ErrFileDoesNotExists = errors.New("file with provided ID does not exists")
)
