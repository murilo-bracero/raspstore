package bootstrap

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/golang-migrate/migrate/v4/source/file"
	"github.com/murilo-bracero/raspstore/file-service/internal/infra/config"
)

type SQLiteBootstraper struct {
}

var _ Bootstraper = (*SQLiteBootstraper)(nil)

// Creates SQLite database file and run migrations against it
func (b *SQLiteBootstraper) Bootstrap(ctx context.Context, config *config.Config) error {
	if ctx.Err() != nil {
		slog.Error("context canceled, aborting bootstrap", "ctxErr", ctx.Err())
		return ctx.Err()
	}

	database, err := sql.Open("sqlite", config.Storage.Path+"/internal/rstore.db")

	if err != nil {
		slog.Error("could not open sqlite database on ./rstore.db", "err", err)
		return err
	}

	dbDriver, err := sqlite.WithInstance(database, &sqlite.Config{})

	if err != nil {
		slog.Error("could initialize sqlite migration", "err", err)
		return err
	}

	fileSource, err := (&file.File{}).Open("file://sql/migrations")

	if err != nil {
		slog.Error("error referencing migrations directory", "err", err)
		return err
	}

	m, err := migrate.NewWithInstance("file", fileSource, "rstore", dbDriver)

	if err != nil {
		slog.Error("error while starting migrations on sqlite database", "err", err)
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		slog.Error("error while executing migrations on sqlite database", "err", err)
		return err
	}

	slog.Info("database migration executed successfully")

	return nil
}
