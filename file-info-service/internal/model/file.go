package model

import (
	"time"
)

type File struct {
	FileId    string `bson:"file_id"`
	Filename  string
	Folder    Folder
	Size      int64
	Owner     string `bson:"owner_user_id"`
	Editors   []string
	Viewers   []string
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
	CreatedBy string    `bson:"created_by"`
	UpdatedBy string    `bson:"updated_by"`
}

type Folder struct {
	Name     string
	IsSecret bool `bson:"is_secret"`
}

type FileMetadataLookup struct {
	FileId    string     `json:"fileId,omitempty" bson:"file_id"`
	Filename  string     `json:"filename,omitempty"`
	Folder    Folder     `json:"folder,omitempty"`
	Size      int64      `json:"size,omitempty"`
	Owner     UserView   `json:"owner,omitempty"`
	Editors   []UserView `json:"editors"`
	Viewers   []UserView `json:"viewers"`
	CreatedAt time.Time  `json:"createdAt,omitempty" bson:"created_at"`
	UpdatedAt time.Time  `json:"updatedAt,omitempty" bson:"updated_at"`
	CreatedBy UserView   `json:"createdBy,omitempty" bson:"created_by"`
	UpdatedBy UserView   `json:"updatedBy,omitempty" bson:"updated_by"`
}

type UserView struct {
	UserId   string `json:"userId,omitempty" bson:"user_id"`
	Username string `json:"username,omitempty"`
}

type FilePage struct {
	Content []*File
	Count   int
}
