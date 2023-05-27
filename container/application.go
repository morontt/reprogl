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

type URLGenerator func(string, bool, ...string) string

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

func RealRemoteAddress(r *http.Request) string {
	addr := r.Header.Get("X-Real-IP")
	if addr == "" {
		addr = r.Header.Get("X-Forwarded-For")
		if addr == "" {
			addr = r.RemoteAddr
		}
	}

	return addr
}

func (app *Application) URLGenerator() URLGenerator {
	return func(routeName string, absoluteURL bool, pairs ...string) string {
		url, err := app.Router.Get(routeName).URL(pairs...)
		if err != nil {
			panic(err)
		}

		var prefix string
		if absoluteURL {
			prefix = "https://" + cnf.Host
		}

		return prefix + url.String()
	}
}
