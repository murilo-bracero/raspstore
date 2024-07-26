package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/murilo-bracero/raspstore/file-service/internal/application/repository"
	"github.com/murilo-bracero/raspstore/file-service/internal/domain/entity"
	"github.com/murilo-bracero/raspstore/file-service/internal/infra/db/gen"
)

type txFilesRepository struct {
	ctx     context.Context
	db      *sql.DB
	queries *gen.Queries
}

var ErrTransactionNotInitialized = errors.New("transaction not initialized for this operation")

var _ repository.TxFilesRepository = (*txFilesRepository)(nil)

func NewTxFilesRepository(ctx context.Context, db *sql.DB) *txFilesRepository {
	return &txFilesRepository{ctx: ctx, queries: gen.New(db), db: db}
}

func (t *txFilesRepository) Begin() (*sql.Tx, error) {
	return t.db.Begin()
}

func (t *txFilesRepository) Commit(tx *sql.Tx) error {
	return tx.Commit()
}

func (t *txFilesRepository) FindById(tx *sql.Tx, userId string, fileId string) (*entity.File, error) {
	nq := t.queries.WithTx(tx)

	rows, err := nq.FindFileByExternalID(t.ctx, gen.FindFileByExternalIDParams{FileID: fileId, OwnerID: userId})

	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return nil, repository.ErrFileDoesNotExists
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

func (t *txFilesRepository) Update(tx *sql.Tx, userId string, file *entity.File) error {
	nq := t.queries.WithTx(tx)

	ts := time.Now()
	file.UpdatedAt = &ts
	file.UpdatedBy = &userId

	return nq.UpdateFileByExternalId(t.ctx, gen.UpdateFileByExternalIdParams{
		FileID:    file.FileId,
		OwnerID:   userId,
		FileName:  file.Filename,
		IsSecret:  file.Secret,
		UpdatedAt: sql.NullInt64{Int64: file.UpdatedAt.UnixMilli(), Valid: true},
		UpdatedBy: sql.NullString{String: *file.UpdatedBy, Valid: true},
	})
}

func (t *txFilesRepository) DeleteFilePermissionByFileId(tx *sql.Tx, fileId string) error {
	nq := t.queries.WithTx(tx)

	return nq.DeleteFilePermissionByFileId(t.ctx, fileId)
}
