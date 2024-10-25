package db

import (
	"database/sql"
	"log/slog"
	"os"

	"github.com/murilo-bracero/raspstore/file-service/internal/infra/config"
	_ "modernc.org/sqlite"
)

type databaseConnection struct {
	db *sql.DB
}

var _ DatabaseConnection = (*databaseConnection)(nil)

func NewSqliteDatabaseConnection(c *config.Config) (*databaseConnection, error) {
	if err := os.MkdirAll(c.Storage.Path+"/internal", os.ModePerm); err != nil {
		return nil, err
	}

	database, err := sql.Open("sqlite", c.Storage.Path+"/internal/rstore.db")

	if err != nil {
		slog.Error("could not open sqlite database on ./rstore.db", "err", err)
		return nil, err
	}

	return &databaseConnection{db: database}, nil
}

func (c *databaseConnection) Db() *sql.DB {
	return c.db
}

func (c *databaseConnection) Close() error {
	return c.db.Close()
}
