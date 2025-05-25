package database

import (
	"database/sql"
	"log"
	"time"
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
	start := time.Now()
	rows, err := dbw.db.Query(query, args...)
	dbw.logQuery(query, args...)
	dbw.logDuration(start)

	return rows, err
}

func (dbw *DB) QueryRow(query string, args ...any) *sql.Row {
	start := time.Now()
	row := dbw.db.QueryRow(query, args...)
	dbw.logQuery(query, args...)
	dbw.logDuration(start)

	return row
}

func (dbw *DB) Exec(query string, args ...any) (sql.Result, error) {
	start := time.Now()
	result, err := dbw.db.Exec(query, args...)
	dbw.logQuery(query, args...)
	dbw.logDuration(start)

	return result, err
}

func (dbw *DB) Close() error {
	return dbw.db.Close()
}

func (dbw *DB) logQuery(query string, args ...any) {
	dbw.logger.Printf("[SQL] query:%s\nparams: %+v\n", query, args)
}

func (dbw *DB) logDuration(t time.Time) {
	dbw.logger.Printf("[SQL] duration %s\n", time.Since(t))
}
