package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"runtime/debug"
	"time"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
	_ "github.com/go-sql-driver/mysql"
	"xelbot.com/reprogl/container"
	"xelbot.com/reprogl/middlewares"
	"xelbot.com/reprogl/views"
)

func main() {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime|log.Lmicroseconds)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	infoLog.Printf("Version: %s (tag: %s)", container.Version, container.GetBuildTag())
	infoLog.Printf("Build time: %s", container.BuildTime)
	infoLog.Printf("Go version: %s", runtime.Version())

	handleError(container.Load("app.ini"), errorLog)

	cfg := container.GetConfig()
	db, err := getDBConnection(cfg.DatabaseDSN, infoLog)
	if err != nil {
		errorLog.Fatal(err)
	}

	goqu.SetTimeLocation(time.Local)

	defer db.Close()

	app := &container.Application{
		ErrorLog: errorLog,
		InfoLog:  infoLog,
		DB:       db,
	}

	router := getRoutes(app)
	handler := middlewares.Session(router, infoLog)
	handler = middlewares.Clacks().Middleware(handler)
	handler = middlewares.CDN(handler)
	handler = middlewares.Recover(handler, app)
	handler = middlewares.Track(handler, app)
	handler = middlewares.AccessLog(handler, app)
	handler = middlewares.ResponseWrapper(handler)

	urlGenerator := func(routeName string, absoluteURL bool, pairs ...string) string {
		url, err := router.Get(routeName).URL(pairs...)
		if err != nil {
			panic(err)
		}

		var prefix string
		if absoluteURL {
			prefix = "https://" + cfg.Host
		}

		return prefix + url.String()
	}

	container.SetURLGenerator(urlGenerator)
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
		logger.Println(fmt.Sprintf("%s", debug.Stack()))
		logger.Fatal(err)
	}
}
