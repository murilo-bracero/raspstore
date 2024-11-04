package repository

import (
	"context"
	"database/sql"

	"github.com/murilo-bracero/raspstore/file-service/internal/domain/entity"
	"github.com/murilo-bracero/raspstore/file-service/internal/infra/db/gen"
)

type usersRepository struct {
	ctx     context.Context
	queries *gen.Queries
}

func NewUsersRepository(ctx context.Context, db *sql.DB) *usersRepository {
	return &usersRepository{ctx: ctx, queries: gen.New(db)}
}

func (r *usersRepository) Save(user *entity.User) error {
	return r.queries.CreateUser(r.ctx, gen.CreateUserParams{
		UserID:       user.Id.String(),
		Username:     user.Username,
		PasswordHash: user.PasswordHash,
		Name:         sql.NullString{String: user.Name, Valid: user.Name != ""},
		RefreshToken: sql.NullString{String: user.RefreshToken, Valid: user.RefreshToken != ""},
	})
}
