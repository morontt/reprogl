package config

import (
	"errors"
	"github.com/gookit/ini/v2"
	"log"
	"os"
)

type AppConfig struct {
	DevMode     bool
	DatabaseDSN string
	HeaderText  string
	Host        string
	Port        string
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
