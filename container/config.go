package container

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/gookit/ini/v2"
)

type AppConfig struct {
	AppData
	CDNBaseURL        string `app_config:"CDN_BASE_URL"`
	DatabaseDSN       string `app_config:"DB"`
	Host              string `app_config:"HOST"`
	Port              string `app_config:"PORT"`
	BackendApiUrl     string `app_config:"BACKEND_API_URL"`
	BackendApiUser    string `app_config:"BACKEND_API_USER"`
	BackendApiWsseKey string `app_config:"BACKEND_API_WSSE_KEY"`
	TelegramToken     string `app_config:"TELEGRAM_TOKEN"`
	TelegramAdminID   int    `app_config:"TELEGRAM_ADMIN_ID"`
	SessionHashKey    string `app_config:"SESSION_HASH_KEY"`
	SessionBlockKey   string `app_config:"SESSION_BLOCK_KEY"`
}

var cnf AppConfig

func Load(file string) error {
	cnf = AppConfig{AppData: loadAppData()}

	return cnf.load(file)
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

func configStringValue(paramName string) (string, error) {
	var value string
	if _, ok := ini.GetValue(paramName); ok {
		value = ini.String(paramName)
	} else {
		return value, errors.New("app.ini: Undefined parameter \"" + paramName + "\"")
	}

	return value, nil
}

func configIntValue(paramName string) (int, error) {
	var value int
	if _, ok := ini.GetValue(paramName); ok {
		value = ini.Int(paramName)
	} else {
		return value, errors.New("app.ini: Undefined parameter \"" + paramName + "\"")
	}

	return value, nil
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

	return setupFields(rv, parseFields(reflect.TypeOf(*c)))
}

func setupFields(rv reflect.Value, fields []structField) error {
	if rv.Kind() != reflect.Pointer {
		return fmt.Errorf("config: argument is not a pointer")
	}

	confValue := rv.Elem()
	for _, field := range fields {
		f := confValue.FieldByName(field.name)
		if f.IsValid() && f.CanSet() {
			switch field.fieldType {
			case reflect.Int:
				intVal, err := configIntValue(field.parameterName)
				if err != nil {
					return err
				}
				x := int64(intVal)
				if !f.OverflowInt(x) {
					f.SetInt(x)
				} else {
					return fmt.Errorf("config: int overflow for field %s", field.name)
				}
			case reflect.String:
				stringVal, err := configStringValue(field.parameterName)
				if err != nil {
					return err
				}
				f.SetString(stringVal)
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
