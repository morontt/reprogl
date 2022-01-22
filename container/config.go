package container

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gookit/ini/v2"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"runtime/debug"
)

type AppConfig struct {
	CDNBaseURL  string
	DevMode     bool
	DatabaseDSN string
	HeaderText  string
	Host        string
	Port        string
}

type Application struct {
	ErrorLog *log.Logger
	InfoLog  *log.Logger
	DB       *sql.DB
	Router   *mux.Router
}

var cnf AppConfig

func init() {
	err := ini.LoadExists("app.ini")
	if err != nil {
		handleError(err)
	}

	if _, ok := ini.GetValue("PORT"); ok {
		cnf.Port = ini.String("PORT")
	} else {
		handleError(errors.New("app.ini: Undefined parameter \"PORT\""))
	}

	if _, ok := ini.GetValue("DEV_MODE"); ok {
		cnf.DevMode = ini.Bool("DEV_MODE")
	} else {
		cnf.DevMode = false
	}

	if _, ok := ini.GetValue("DB"); ok {
		cnf.DatabaseDSN = ini.String("DB")
	} else {
		handleError(errors.New("app.ini: Undefined parameter \"DB\""))
	}

	if _, ok := ini.GetValue("HOST"); ok {
		cnf.Host = ini.String("HOST")
	} else {
		handleError(errors.New("app.ini: Undefined parameter \"HOST\""))
	}

	if _, ok := ini.GetValue("HEADER_TEXT"); ok {
		cnf.HeaderText = ini.String("HEADER_TEXT")
	} else {
		handleError(errors.New("app.ini: Undefined parameter \"HEADER_TEXT\""))
	}

	if _, ok := ini.GetValue("CDN_BASE_URL"); ok {
		cnf.CDNBaseURL = ini.String("CDN_BASE_URL")
	} else {
		handleError(errors.New("app.ini: Undefined parameter \"CDN_BASE_URL\""))
	}
}

func Get() AppConfig {
	return cnf
}

func IsDevMode() bool {
	return cnf.DevMode
}

func handleError(err error) {
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	errorLog.Fatal(err)
}

func (app *Application) NotFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *Application) ServerError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Println(trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *Application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}
