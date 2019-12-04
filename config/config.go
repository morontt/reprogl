package config

import (
	"errors"
	ini "github.com/gookit/ini/v2"
)

type AppConfig struct {
	Port    string
	DevMode bool
}

var cnf AppConfig

func Load() error {
	err := ini.LoadExists("app.ini")
	if err != nil {
		return err
	}

	if _, ok := ini.GetValue("PORT"); ok {
		cnf.Port = ini.String("PORT")
	} else {
		return errors.New("app.ini: Undefined parameter \"PORT\"")
	}

	if _, ok := ini.GetValue("DEV_MODE"); ok {
		cnf.DevMode = ini.Bool("DEV_MODE")
	} else {
		return errors.New("app.ini: Undefined parameter \"DEV_MODE\"")
	}

	return nil
}

func Get() AppConfig {
	return cnf
}
