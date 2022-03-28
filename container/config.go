package container

import (
	"errors"
	"github.com/gookit/ini/v2"
	"log"
	"os"
)

// TODO Auto read config by parameter type and name

type AppConfig struct {
	CDNBaseURL  string
	DevMode     bool
	DatabaseDSN string
	HeaderText  string
	Host        string
	Port        string
	Author      string
	AdminEmail  string
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

	if _, ok := ini.GetValue("AUTHOR"); ok {
		cnf.Author = ini.String("AUTHOR")
	} else {
		handleError(errors.New("app.ini: Undefined parameter \"AUTHOR\""))
	}

	if _, ok := ini.GetValue("ADMIN_EMAIL"); ok {
		cnf.AdminEmail = ini.String("ADMIN_EMAIL")
	} else {
		handleError(errors.New("app.ini: Undefined parameter \"ADMIN_EMAIL\""))
	}
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
