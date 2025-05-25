package database

import (
	"database/sql"
	"log"
)

type DB struct {
	db     *sql.DB
	logger *log.Logger
}

func Wrap(db *sql.DB, logger *log.Logger) *DB {
	return &DB{
		db:     db,
		logger: logger,
	}
}

func (dbw *DB) Query(query string, args ...any) (*sql.Rows, error) {
	return dbw.db.Query(query, args...)
}

func (dbw *DB) QueryRow(query string, args ...any) *sql.Row {
	return dbw.db.QueryRow(query, args...)
}

func (dbw *DB) Exec(query string, args ...any) (sql.Result, error) {
	return dbw.db.Exec(query, args...)
}

func (dbw *DB) Close() error {
	return dbw.db.Close()
}
