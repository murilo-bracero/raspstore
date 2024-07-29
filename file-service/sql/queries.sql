-- name: FindFileByExternalID :many
SELECT f.*, fp.*
FROM files f
LEFT JOIN files_permissions fp ON f.file_id = fp.file_id
WHERE f.file_id = ?1
AND (
    f.owner_id = ?2 OR EXISTS (
        SELECT 1
        FROM files_permissions ffp
        WHERE ffp.file_id = f.file_id AND ffp.user_id = ?2
    )
);

-- name: CreateFile :exec
INSERT INTO files (file_id, file_name, size, is_secret, owner_id, created_at, created_by)
VALUES (?, ?, ?, ?, ?, ?, ?);

-- name: DeleteFileByExternalID :exec
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
);

-- name: UpdateFileByExternalId :exec
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
);

-- name: DeleteFilePermissionByFileId :exec
DELETE FROM files_permissions WHERE file_id = ?;

-- name: FindAllFiles :many
SELECT f.*, COUNT() OVER() AS totalCount
FROM files f
LEFT JOIN files_permissions fp ON f.file_id = fp.file_id
WHERE (f.owner_id = ?1 OR fp.user_id = ?1)
AND f.file_name LIKE ?2
AND f.is_secret = ?3
ORDER BY f.created_at DESC
LIMIT ?4
OFFSET ?5;

-- name: FindUsageByUserId :one
SELECT SUM(f.size) as totalSize
FROM files f
WHERE f.owner_id = ?
GROUP BY f.owner_id;