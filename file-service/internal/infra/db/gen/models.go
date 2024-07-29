// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package gen

import (
	"database/sql"
)

type File struct {
	FileID    string
	FileName  string
	Size      int64
	IsSecret  bool
	OwnerID   string
	CreatedAt int64
	UpdatedAt sql.NullInt64
	CreatedBy string
	UpdatedBy sql.NullString
}

type FilesPermission struct {
	PermissionID string
	FileID       string
	Permission   string
	UserID       string
}