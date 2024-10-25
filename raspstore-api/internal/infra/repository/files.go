package repository

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/murilo-bracero/raspstore/file-service/internal/domain/entity"
	"github.com/murilo-bracero/raspstore/file-service/internal/infra/db/gen"
)

type filesRepository struct {
	ctx     context.Context
	queries *gen.Queries
}

var _ FilesRepository = (*filesRepository)(nil)

func NewFilesRepository(ctx context.Context, db *sql.DB) *filesRepository {
	return &filesRepository{queries: gen.New(db), ctx: ctx}
}

func (r *filesRepository) Save(file *entity.File) error {
	file.FileId = uuid.NewString()
	file.CreatedAt = time.Now()
	ts := time.Now()
	file.UpdatedAt = &ts

	err := r.queries.CreateFile(r.ctx, gen.CreateFileParams{
		FileName:  file.Filename,
		Size:      file.Size,
		IsSecret:  file.Secret,
		OwnerID:   file.Owner,
		FileID:    file.FileId,
		CreatedAt: file.CreatedAt.UnixMilli(),
		CreatedBy: file.Owner,
	})

	if err != nil {
		slog.Error("could not save file into database", "err", err)
		return err
	}

	return nil
}

func (r *filesRepository) FindById(id string, userId string) (*entity.File, error) {
	rows, err := r.queries.FindFileByID(r.ctx, gen.FindFileByIDParams{FileID: id, OwnerID: userId})

	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return nil, ErrFileDoesNotExists
	}

	ref := rows[0]

	updatedAt := time.UnixMilli(ref.UpdatedAt.Int64)

	file := &entity.File{
		FileId:    ref.FileID,
		Filename:  ref.FileName,
		Size:      ref.Size,
		Secret:    ref.IsSecret,
		Owner:     ref.OwnerID,
		CreatedAt: time.UnixMilli(ref.CreatedAt),
		UpdatedAt: &updatedAt,
		CreatedBy: ref.CreatedBy,
		UpdatedBy: &ref.UpdatedBy.String,
		Viewers:   []string{},
		Editors:   []string{},
	}

	for _, row := range rows {
		if !row.Permission.Valid {
			continue
		}

		permission := row.Permission.String

		if permission == "VIEWER" {
			file.Viewers = append(file.Viewers, row.UserID.String)
		}

		if permission == "EDITOR" {
			file.Editors = append(file.Editors, row.UserID.String)
		}
	}

	return file, nil
}

func (r *filesRepository) Delete(userId string, fileId string) error {
	return r.queries.DeleteFileByID(r.ctx, gen.DeleteFileByIDParams{FileID: fileId, OwnerID: userId})
}

func (r *filesRepository) Update(userId string, file *entity.File) error {
	ts := time.Now()
	file.UpdatedAt = &ts
	file.UpdatedBy = &userId

	return r.queries.UpdateFileByID(r.ctx, gen.UpdateFileByIDParams{
		FileID:    file.FileId,
		OwnerID:   userId,
		FileName:  file.Filename,
		IsSecret:  file.Secret,
		UpdatedAt: sql.NullInt64{Int64: file.UpdatedAt.UnixMilli(), Valid: true},
		UpdatedBy: sql.NullString{String: *file.UpdatedBy, Valid: true},
	})
}

func (r *filesRepository) FindAll(userId string, page int, size int, filename string, secret bool) (filesPage *entity.FilePage, err error) {
	rows, err := r.queries.FindAllFiles(r.ctx, gen.FindAllFilesParams{
		OwnerID:  userId,
		FileName: "%" + filename + "%",
		IsSecret: secret,
		Limit:    int64(size),
		Offset:   int64(page) * int64(size),
	})

	if err != nil {
		return nil, err
	}

	totalCount := 0

	if len(rows) != 0 {
		totalCount = int(rows[0].Totalcount)
	}

	filePage := &entity.FilePage{Count: totalCount, Content: make([]*entity.File, len(rows))}

	for i, row := range rows {
		filePage.Content[i] = &entity.File{
			FileId:    row.FileID,
			Filename:  row.FileName,
			Size:      row.Size,
			Secret:    row.IsSecret,
			Owner:     row.OwnerID,
			CreatedAt: time.UnixMilli(row.CreatedAt),
			CreatedBy: row.CreatedBy,
		}

		if row.UpdatedAt.Valid {
			updatedAt := time.UnixMilli(row.UpdatedAt.Int64)
			filePage.Content[i].UpdatedAt = &updatedAt
		}

		if row.UpdatedBy.Valid {
			updatedBy := row.UpdatedBy.String
			filePage.Content[i].UpdatedBy = &updatedBy
		}
	}

	return filePage, nil
}

func (r *filesRepository) FindUsageByUserId(userId string) (int64, error) {
	row, err := r.queries.FindUsageByUserID(r.ctx, userId)

	if err == sql.ErrNoRows {
		return 0, nil
	}

	if err != nil {
		return 0, err
	}

	return int64(row.Float64), nil
}

func (r *filesRepository) DeleteFilePermissionByFileId(fileId string) error {
	return r.queries.DeleteFilePermissionByFileID(r.ctx, fileId)
}
