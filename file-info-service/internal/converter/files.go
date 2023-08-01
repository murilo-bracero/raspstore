package converter

import (
	"time"

	"github.com/google/uuid"
	"github.com/murilo-bracero/raspstore/file-service/internal/model"
)

func CreateFile(filename string, size int64, secret bool, ownerId string) *model.File {
	return &model.File{
		FileId:    uuid.NewString(),
		Filename:  filename,
		Size:      size,
		Secret:    secret,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Viewers:   []string{},
		Editors:   []string{},
		CreatedBy: ownerId,
		UpdatedBy: ownerId,
		Owner:     ownerId,
	}
}
