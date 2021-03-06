package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"os"
	"time"
	"xelbot.com/reprogl/config"
	"xelbot.com/reprogl/middlewares"
	"xelbot.com/reprogl/views"
)

func main() {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime|log.Lmicroseconds)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	cfg := config.Get()
	infoLog.Print("Trying to connect to the database")
	db, err := openDB(cfg.DatabaseDSN)
	if err != nil {
		errorLog.Fatal(err)
	}

	app := &config.Application{
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

func openDB(dsn string) (*sql.DB, error) {
	time.Sleep(500 * time.Millisecond)
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
