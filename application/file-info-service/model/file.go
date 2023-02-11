package model

import (
	"time"
)

type File struct {
	FileId    string    `bson:"file_id"`
	Filename  string    `bson:"filename"`
	Path      string    `bson:"path"`
	Size      int64     `bson:"size"`
	Owner     string    `bson:"owner_user_id"`
	Editors   []string  `bson:"editors"`
	Viewers   []string  `bson:"viewers"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
	CreatedBy string    `bson:"created_by"`
	UpdatedBy string    `bson:"updated_by"`
}

type FileMetadataLookup struct {
	FileId    string     `json:"fileId,omitempty" bson:"file_id"`
	Filename  string     `json:"path,omitempty" bson:"filename"`
	Path      string     `json:"filename,omitempty" bson:"path"`
	Size      int64      `json:"size,omitempty" bson:"size"`
	Owner     string     `json:"owner,omitempty" bson:"owner_user_id"`
	Editors   []UserView `json:"editors" bson:"editors"`
	Viewers   []UserView `json:"viewers" bson:"viewers"`
	CreatedAt time.Time  `json:"createdBy,omitempty" bson:"created_at"`
	UpdatedAt time.Time  `json:"updatedBy,omitempty" bson:"updated_at"`
	CreatedBy UserView   `json:"createdAt,omitempty" bson:"created_by"`
	UpdatedBy UserView   `json:"updatedAt,omitempty" bson:"updated_by"`
}

type UserView struct {
	UserId   string `json:"userId,omitempty" bson:"user_id"`
	Username string `json:"username,omitempty" bson:"username"`
}

type FilePage struct {
	Content []*File `bson:"content"`
	Count   int     `bson:"count"`
}
