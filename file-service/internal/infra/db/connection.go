package db

import "database/sql"

type DatabaseConnection interface {
	Db() *sql.DB
	Close() error
}
