package db

import (
	"database/sql"
	"log/slog"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
	"github.com/murilo-bracero/raspstore/file-service/internal/infra/config"
)

type databaseConnection struct {
	db *sql.DB
}

var _ DatabaseConnection = (*databaseConnection)(nil)

func NewSqliteDatabaseConnection(c *config.Config) (*databaseConnection, error) {
	if err := os.MkdirAll(c.Server.Storage.Path+"/internal", os.ModePerm); err != nil {
		return nil, err
	}

	slog.Info(c.Server.Storage.Path + "/internal/rstore.db")

	database, err := sql.Open("sqlite3", c.Server.Storage.Path+"/internal/rstore.db")

	if err != nil {
		slog.Error("could not open sqlite database on ./rstore.db", "err", err)
		return nil, err
	}

	dbDriver, err := sqlite3.WithInstance(database, &sqlite3.Config{})

	if err != nil {
		slog.Error("could initialize sqlite migration", "err", err)
		return nil, err
	}

	fileSource, err := (&file.File{}).Open("file://db/migrations")

	if err != nil {
		slog.Error("error referencing migrations directory", "err", err)
		return nil, err
	}

	m, err := migrate.NewWithInstance("file", fileSource, "rstore", dbDriver)

	if err != nil {
		slog.Error("error while starting migrations on sqlite database", "err", err)
		return nil, err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		slog.Error("error while executing migrations on sqlite database", "err", err)
		return nil, err
	}

	slog.Info("database migration executed successfully")

	return &databaseConnection{db: database}, nil
}

func (c *databaseConnection) Db() *sql.DB {
	return c.db
}

func (c *databaseConnection) Close() error {
	return c.db.Close()
}
