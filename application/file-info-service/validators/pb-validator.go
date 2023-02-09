package validators

import (
	"errors"
)

var (
	ErrReceiveStreaming = errors.New("an error occured while receiving the stream, try again later")
	ErrUploadFile       = errors.New("an error occured while writing uploaded file in server. try again later")
	ErrCreatedByEmpty   = errors.New("filed CreatedBy must not be empty")
	ErrFilenameEmpty    = errors.New("field Filename must not be empty")
	ErrWrongID          = errors.New("provided id is invalid")
)
