package container

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"runtime"
	"runtime/debug"
	"time"

	"github.com/xelbot/yetacache"
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

	intCache *yetacache.Cache[string, int]
	strCache *yetacache.Cache[string, string]
}

var urlGen URLGenerator

func init() {
	re := regexp.MustCompile(`^\D*(\d+\.\d+(?:\.\d+)?)`)
	GoVersionNumbers = re.FindStringSubmatch(runtime.Version())[1]

	checkBuildFlags()
}

func checkBuildFlags() {
	if len(Version) == 0 {
		panic("ldflags: xelbot.com/reprogl/container.Version is empty")
	}
	if len(GitRevision) == 0 {
		panic("ldflags: xelbot.com/reprogl/container.GitRevision is empty")
	}

	_, err := time.Parse(time.RFC1123, BuildTime)
	if err != nil {
		panic("ldflags: xelbot.com/reprogl/container.BuildTime wrong format")
	}
}

func SetURLGenerator(u URLGenerator) {
	urlGen = u
}

func GenerateURL(routeName string, pairs ...string) string {
	return urlGen(routeName, false, pairs...)
}

func GenerateAbsoluteURL(routeName string, pairs ...string) string {
	return urlGen(routeName, true, pairs...)
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

func (app *Application) LogError(err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Println(trace)
}

func (app *Application) GetIntCache() *yetacache.Cache[string, int] {
	if app.intCache == nil {
		app.InfoLog.Println("[CACHE] create integer instance")
		app.intCache = yetacache.New[string, int](time.Hour, 8*time.Hour)
	}

	return app.intCache
}

func (app *Application) GetStringCache() *yetacache.Cache[string, string] {
	if app.strCache == nil {
		app.InfoLog.Println("[CACHE] create string instance")
		app.strCache = yetacache.New[string, string](time.Hour, 8*time.Hour)
	}

	return app.strCache
}

func (app *Application) Stop() error {
	err := app.DB.Close()
	if err != nil {
		return err
	}
	app.InfoLog.Print("The database connection is closed")

	return nil
}
