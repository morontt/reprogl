package container

import (
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"

	"github.com/gookit/ini/v2"
)

type AppConfig struct {
	CDNBaseURL        string `app_config:"CDN_BASE_URL"`
	DatabaseDSN       string `app_config:"DB"`
	HeaderText        string `app_config:"HEADER_TEXT"`
	Host              string `app_config:"HOST"`
	Port              string `app_config:"PORT"`
	Author            string `app_config:"AUTHOR"`
	AuthorBio         string `app_config:"AUTHOR_BIO"`
	AuthorGithub      string `app_config:"AUTHOR_GITHUB"`
	AuthorTelegram    string `app_config:"AUTHOR_TELEGRAM"`
	AdminEmail        string `app_config:"ADMIN_EMAIL"`
	BackendApiUrl     string `app_config:"BACKEND_API_URL"`
	BackendApiUser    string `app_config:"BACKEND_API_USER"`
	BackendApiWsseKey string `app_config:"BACKEND_API_WSSE_KEY"`
	TelegramToken     string `app_config:"TELEGRAM_TOKEN"`
	TelegramAdminID   int    `app_config:"TELEGRAM_ADMIN_ID"`
}

var cnf AppConfig

func init() {
	err := cnf.load("app.ini")
	if err != nil {
		handleError(err)
	}
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

type structField struct {
	name          string
	parameterName string
	fieldType     reflect.Kind
}

func (c *AppConfig) load(file string) error {
	err := ini.LoadExists(file)
	if err != nil {
		return err
	}

	rv := reflect.ValueOf(c)
	if rv.IsNil() {
		return fmt.Errorf("config: load error %v", reflect.TypeOf(c))
	}

	fields := parseFields(reflect.TypeOf(*c))

	confValue := rv.Elem()
	for _, field := range fields {
		f := confValue.FieldByName(field.name)
		if f.IsValid() && f.CanSet() {
			switch field.fieldType {
			case reflect.Int:
				x := int64(configIntValue(field.parameterName))
				if !f.OverflowInt(x) {
					f.SetInt(x)
				} else {
					return fmt.Errorf("config: int overflow for field %s", field.name)
				}
			case reflect.String:
				f.SetString(configStringValue(field.parameterName))
			default:
				return fmt.Errorf("config: undefined type %s of field %s", field.fieldType, field.name)
			}
		}
	}

	return nil
}

func parseFields(t reflect.Type) []structField {
	var fields = make([]structField, 0)
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		tag := f.Tag.Get("app_config")
		if tag != "" {
			fields = append(fields, structField{
				parameterName: tag,
				fieldType:     f.Type.Kind(),
				name:          f.Name,
			})
		}
	}

	return fields
}
