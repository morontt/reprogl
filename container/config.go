package container

import (
	"errors"
	"github.com/gookit/ini/v2"
	"log"
	"os"
)

// TODO Auto read config by parameter type and name

type AppConfig struct {
	CDNBaseURL        string
	DevMode           bool
	DatabaseDSN       string
	HeaderText        string
	Host              string
	Port              string
	Author            string
	AuthorBio         string
	AuthorGithub      string
	AuthorInsta       string
	AdminEmail        string
	BackendApiUrl     string
	BackendApiUser    string
	BackendApiWsseKey string
}

var cnf AppConfig

func init() {
	err := ini.LoadExists("app.ini")
	if err != nil {
		handleError(err)
	}

	if _, ok := ini.GetValue("DEV_MODE"); ok {
		cnf.DevMode = ini.Bool("DEV_MODE")
	} else {
		cnf.DevMode = false
	}

	cnf.Port = configStringValue("PORT")
	cnf.DatabaseDSN = configStringValue("DB")
	cnf.Host = configStringValue("HOST")
	cnf.HeaderText = configStringValue("HEADER_TEXT")
	cnf.CDNBaseURL = configStringValue("CDN_BASE_URL")
	cnf.Author = configStringValue("AUTHOR")
	cnf.AdminEmail = configStringValue("ADMIN_EMAIL")
	cnf.AuthorBio = configStringValue("AUTHOR_BIO")
	cnf.AuthorGithub = configStringValue("AUTHOR_GITHUB")
	cnf.AuthorInsta = configStringValue("AUTHOR_INSTAGRAM")
	cnf.BackendApiUrl = configStringValue("BACKEND_API_URL")
	cnf.BackendApiUser = configStringValue("BACKEND_API_USER")
	cnf.BackendApiWsseKey = configStringValue("BACKEND_API_WSSE_KEY")
}

func GetConfig() AppConfig {
	return cnf
}

func IsDevMode() bool {
	return cnf.DevMode
}

func handleError(err error) {
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	errorLog.Fatal(err)
}

func configStringValue(paramName string) (value string) {
	if _, ok := ini.GetValue(paramName); ok {
		value = ini.String(paramName)
	} else {
		handleError(errors.New("app.ini: Undefined parameter \"" + paramName + "\""))
	}

	return
}
