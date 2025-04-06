package container

import (
	"database/sql"
	"log"
	"time"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
	_ "github.com/go-sql-driver/mysql"
)

func (app *Application) SetupDatabase() error {
	db, err := getDBConnection(app.InfoLog)
	if err != nil {
		return err
	}

	app.DB = db
	goqu.SetTimeLocation(time.Local)

	return nil
}

func getDBConnection(logger *log.Logger) (db *sql.DB, err error) {
	var i int

	for i < 5 {
		logger.Print("Trying to connect to the database")
		db, err = openDB(cnf.DatabaseDSN)
		if err == nil {
			logger.Print("The database is connected")

			return
		} else {
			logger.Print(err)
		}

		i++
		time.Sleep(1000 * time.Millisecond)
	}

	return nil, err
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
