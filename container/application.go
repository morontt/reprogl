package container

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"regexp"
	"runtime"
	"runtime/debug"
)

var Version string
var GitRevision string
var BuildTime string
var GoVersionNumbers string

type Application struct {
	ErrorLog *log.Logger
	InfoLog  *log.Logger
	DB       *sql.DB
	Router   *mux.Router
}

func init() {
	re := regexp.MustCompile(`^\D*(\d+\.\d+(?:\.\d+)?)`)
	GoVersionNumbers = re.FindStringSubmatch(runtime.Version())[1]
}

func (app *Application) NotFound(w http.ResponseWriter) {
	app.ClientError(w, http.StatusNotFound)
}

func (app *Application) ServerError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Println(trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *Application) ClientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}
