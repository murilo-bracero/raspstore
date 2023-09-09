package model

import (
	"time"
)

type File struct {
	FileId    string    `json:"fileId,omitempty" bson:"file_id"`
	Filename  string    `json:"filename,omitempty"`
	Size      int64     `json:"size,omitempty"`
	Secret    bool      `json:"-" bson:"is_secret"`
	Owner     string    `json:"owner,omitempty"`
	Editors   []string  `json:"editors"`
	Viewers   []string  `json:"viewers"`
	CreatedAt time.Time `json:"createdAt,omitempty" bson:"created_at"`
	UpdatedAt time.Time `json:"updatedAt,omitempty" bson:"updated_at"`
	CreatedBy string    `json:"createdBy,omitempty" bson:"created_by"`
	UpdatedBy string    `json:"updatedBy,omitempty" bson:"updated_by"`
}

type FilePage struct {
	Content []*File
	Count   int
}
