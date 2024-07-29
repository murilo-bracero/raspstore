// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: queries.sql

package gen

import (
	"context"
	"database/sql"
)

const createFile = `-- name: CreateFile :exec
INSERT INTO files (file_id, file_name, size, is_secret, owner_id, created_at, created_by)
VALUES (?, ?, ?, ?, ?, ?, ?)
`

type CreateFileParams struct {
	FileID    string
	FileName  string
	Size      int64
	IsSecret  bool
	OwnerID   string
	CreatedAt int64
	CreatedBy string
}

func (q *Queries) CreateFile(ctx context.Context, arg CreateFileParams) error {
	_, err := q.db.ExecContext(ctx, createFile,
		arg.FileID,
		arg.FileName,
		arg.Size,
		arg.IsSecret,
		arg.OwnerID,
		arg.CreatedAt,
		arg.CreatedBy,
	)
	return err
}

const deleteFileByExternalID = `-- name: DeleteFileByExternalID :exec
DELETE FROM files
WHERE file_id IN (
    SELECT f.file_id 
    FROM files f
    LEFT JOIN files_permissions fp ON f.file_id = fp.file_id AND fp.permission = 'EDITOR' 
    WHERE f.file_id = ?1
    AND (
        f.owner_id = ?2 OR
        fp.user_id = ?2
    )
)
`

type DeleteFileByExternalIDParams struct {
	FileID  string
	OwnerID string
}

func (q *Queries) DeleteFileByExternalID(ctx context.Context, arg DeleteFileByExternalIDParams) error {
	_, err := q.db.ExecContext(ctx, deleteFileByExternalID, arg.FileID, arg.OwnerID)
	return err
}

const deleteFilePermissionByFileId = `-- name: DeleteFilePermissionByFileId :exec
DELETE FROM files_permissions WHERE file_id = ?
`

func (q *Queries) DeleteFilePermissionByFileId(ctx context.Context, fileID string) error {
	_, err := q.db.ExecContext(ctx, deleteFilePermissionByFileId, fileID)
	return err
}

const findAllFiles = `-- name: FindAllFiles :many
SELECT f.file_id, f.file_name, f.size, f.is_secret, f.owner_id, f.created_at, f.updated_at, f.created_by, f.updated_by, COUNT() OVER() AS totalCount
FROM files f
LEFT JOIN files_permissions fp ON f.file_id = fp.file_id
WHERE (f.owner_id = ?1 OR fp.user_id = ?1)
AND f.file_name LIKE ?2
AND f.is_secret = ?3
ORDER BY f.created_at DESC
LIMIT ?4
OFFSET ?5
`

type FindAllFilesParams struct {
	OwnerID  string
	FileName string
	IsSecret bool
	Limit    int64
	Offset   int64
}

type FindAllFilesRow struct {
	FileID     string
	FileName   string
	Size       int64
	IsSecret   bool
	OwnerID    string
	CreatedAt  int64
	UpdatedAt  sql.NullInt64
	CreatedBy  string
	UpdatedBy  sql.NullString
	Totalcount int64
}

func (q *Queries) FindAllFiles(ctx context.Context, arg FindAllFilesParams) ([]FindAllFilesRow, error) {
	rows, err := q.db.QueryContext(ctx, findAllFiles,
		arg.OwnerID,
		arg.FileName,
		arg.IsSecret,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FindAllFilesRow
	for rows.Next() {
		var i FindAllFilesRow
		if err := rows.Scan(
			&i.FileID,
			&i.FileName,
			&i.Size,
			&i.IsSecret,
			&i.OwnerID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.CreatedBy,
			&i.UpdatedBy,
			&i.Totalcount,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const findFileByExternalID = `-- name: FindFileByExternalID :many
SELECT f.file_id, f.file_name, f.size, f.is_secret, f.owner_id, f.created_at, f.updated_at, f.created_by, f.updated_by, fp.permission_id, fp.file_id, fp.permission, fp.user_id
FROM files f
LEFT JOIN files_permissions fp ON f.file_id = fp.file_id
WHERE f.file_id = ?1
AND (
    f.owner_id = ?2 OR EXISTS (
        SELECT 1
        FROM files_permissions ffp
        WHERE ffp.file_id = f.file_id AND ffp.user_id = ?2
    )
)
`

type FindFileByExternalIDParams struct {
	FileID  string
	OwnerID string
}

type FindFileByExternalIDRow struct {
	FileID       string
	FileName     string
	Size         int64
	IsSecret     bool
	OwnerID      string
	CreatedAt    int64
	UpdatedAt    sql.NullInt64
	CreatedBy    string
	UpdatedBy    sql.NullString
	PermissionID sql.NullString
	FileID_2     sql.NullString
	Permission   sql.NullString
	UserID       sql.NullString
}

func (q *Queries) FindFileByExternalID(ctx context.Context, arg FindFileByExternalIDParams) ([]FindFileByExternalIDRow, error) {
	rows, err := q.db.QueryContext(ctx, findFileByExternalID, arg.FileID, arg.OwnerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FindFileByExternalIDRow
	for rows.Next() {
		var i FindFileByExternalIDRow
		if err := rows.Scan(
			&i.FileID,
			&i.FileName,
			&i.Size,
			&i.IsSecret,
			&i.OwnerID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.CreatedBy,
			&i.UpdatedBy,
			&i.PermissionID,
			&i.FileID_2,
			&i.Permission,
			&i.UserID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const findUsageByUserId = `-- name: FindUsageByUserId :one
SELECT SUM(f.size) as totalSize
FROM files f
WHERE f.owner_id = ?
GROUP BY f.owner_id
`

func (q *Queries) FindUsageByUserId(ctx context.Context, ownerID string) (sql.NullFloat64, error) {
	row := q.db.QueryRowContext(ctx, findUsageByUserId, ownerID)
	var totalsize sql.NullFloat64
	err := row.Scan(&totalsize)
	return totalsize, err
}

const updateFileByExternalId = `-- name: UpdateFileByExternalId :exec
UPDATE files SET 
file_name = ?3,
is_secret = ?4,
updated_at = ?5,
updated_by = ?6
WHERE file_id IN (
    SELECT f.file_id
    FROM files f
    LEFT JOIN files_permissions fp ON f.file_id = fp.file_id AND fp.permission = 'EDITOR' 
    WHERE f.file_id = ?1
    AND (
        f.owner_id = ?2 OR
        fp.user_id = ?2
    )
)
`

type UpdateFileByExternalIdParams struct {
	FileID    string
	OwnerID   string
	FileName  string
	IsSecret  bool
	UpdatedAt sql.NullInt64
	UpdatedBy sql.NullString
}

func (q *Queries) UpdateFileByExternalId(ctx context.Context, arg UpdateFileByExternalIdParams) error {
	_, err := q.db.ExecContext(ctx, updateFileByExternalId,
		arg.FileID,
		arg.OwnerID,
		arg.FileName,
		arg.IsSecret,
		arg.UpdatedAt,
		arg.UpdatedBy,
	)
	return err
}