package container

import (
	"errors"
	"log"
	"os"

	"github.com/gookit/ini/v2"
)

// TODO Auto read config by parameter type and name

type AppConfig struct {
	CDNBaseURL        string
	DatabaseDSN       string
	HeaderText        string
	Host              string
	Port              string
	Author            string
	AuthorBio         string
	AuthorGithub      string
	AuthorTelegram    string
	AdminEmail        string
	BackendApiUrl     string
	BackendApiUser    string
	BackendApiWsseKey string
	TelegramToken     string
	TelegramAdminID   int
}

var cnf AppConfig

func init() {
	err := ini.LoadExists("app.ini")
	if err != nil {
		handleError(err)
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
	cnf.AuthorTelegram = configStringValue("AUTHOR_TELEGRAM")
	cnf.BackendApiUrl = configStringValue("BACKEND_API_URL")
	cnf.BackendApiUser = configStringValue("BACKEND_API_USER")
	cnf.BackendApiWsseKey = configStringValue("BACKEND_API_WSSE_KEY")
	cnf.TelegramToken = configStringValue("TELEGRAM_TOKEN")
	cnf.TelegramAdminID = configIntValue("TELEGRAM_ADMIN_ID")
}

func GetConfig() AppConfig {
	return cnf
}

func IsDevMode() bool {
	return devMode
}

func GetBuildTag() (tag string) {
	if devMode {
		tag = "dev"
	} else {
		tag = "prod"
	}

	return
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

func configIntValue(paramName string) (value int) {
	if _, ok := ini.GetValue(paramName); ok {
		value = ini.Int(paramName)
	} else {
		handleError(errors.New("app.ini: Undefined parameter \"" + paramName + "\""))
	}

	return
}
