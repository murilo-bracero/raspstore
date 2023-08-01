package internal

import "errors"

var (
	ErrReceiveStreaming  = errors.New("an error occured while receiving the stream, try again later")
	ErrUploadFile        = errors.New("an error occured while writing uploaded file in server. try again later")
	ErrOwnerIdEmpty      = errors.New("filed OwnerId must not be empty")
	ErrFilenameEmpty     = errors.New("field Filename must not be empty")
	ErrWrongID           = errors.New("provided id is invalid")
	ErrEmptyFileId       = errors.New("provided id is empty")
	ErrPathEmpty         = errors.New("provided field path must not be null or empty")
	ErrInvalidSize       = errors.New("provided field size must be greatter than 0")
	ErrEditorsNil        = errors.New("editors field must be an array")
	ErrViewersNil        = errors.New("viewers field must be an array")
	ErrNotAvailableSpace = errors.New("file is greather than the space available for your user")
	ErrFileDoesNotExists = errors.New("file with provided ID does not exists")
)