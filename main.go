package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"os"
	"time"
	"xelbot.com/reprogl/container"
	"xelbot.com/reprogl/middlewares"
	"xelbot.com/reprogl/views"
)

func main() {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime|log.Lmicroseconds)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	cfg := container.GetConfig()
	db, err := getDBConnection(cfg.DatabaseDSN, infoLog)
	if err != nil {
		errorLog.Fatal(err)
	}

	app := &container.Application{
		ErrorLog: errorLog,
		InfoLog:  infoLog,
		DB:       db,
	}

	router := getRoutes(app)
	handler := middlewares.Recover(router, app)
	handler = middlewares.AccessLog(handler, app)

	views.SetRouter(router)
	handleError(views.LoadViewSet(), errorLog)

	server := &http.Server{
		Handler:      handler,
		Addr:         ":" + cfg.Port,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}

	infoLog.Printf("Starting server on %s port", cfg.Port)
	handleError(server.ListenAndServe(), errorLog)
}

func getDBConnection(dsn string, logger *log.Logger) (db *sql.DB, err error) {
	var i int

	for i < 5 {
		logger.Print("Trying to connect to the database")
		db, err = openDB(dsn)
		if err == nil {
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

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func handleError(err error, logger *log.Logger) {
	if err != nil {
		logger.Fatal(err)
	}
}
