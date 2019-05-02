package config

import (
	"errors"
	"github.com/gookit/ini"
)

type AppConfig struct {
	Port string
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

	return nil
}

func Get() AppConfig {
	return cnf
}
