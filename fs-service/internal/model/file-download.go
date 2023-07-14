package model

import "os"

type FileDownloadRepresentation struct {
	Filename string
	File     *os.File
	FileSize int64
}
