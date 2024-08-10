package entity

import (
	"time"

	"github.com/google/uuid"
)

type File struct {
	FileId    string     `json:"fileId,omitempty" bson:"file_id"`
	Filename  string     `json:"filename,omitempty"`
	Size      int64      `json:"size,omitempty"`
	Secret    bool       `json:"secret" bson:"is_secret"`
	Owner     string     `json:"owner,omitempty"`
	Editors   []string   `json:"editors"`
	Viewers   []string   `json:"viewers"`
	CreatedAt time.Time  `json:"createdAt,omitempty" bson:"created_at"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty" bson:"updated_at"`
	CreatedBy string     `json:"createdBy,omitempty" bson:"created_by"`
	UpdatedBy *string    `json:"updatedBy,omitempty" bson:"updated_by"`
}

func NewFile(filename string, size int64, secret bool, ownerId string) *File {
	return &File{
		FileId:    uuid.NewString(),
		Filename:  filename,
		Size:      size,
		Secret:    secret,
		CreatedAt: time.Now(),
		Viewers:   []string{},
		Editors:   []string{},
		CreatedBy: ownerId,
		Owner:     ownerId,
	}
}

type FilePage struct {
	Content []*File
	Count   int
}
