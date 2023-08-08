package configs

import (
	"fmt"
	"reflect"

	env "github.com/caarlos0/env/v6"
	"github.com/timdevlet/todo/internal/helpers"
)

func NewConfigsFromEnv() *Configs {
	o := Configs{}

	return o.loadFromEnv()
}

type Configs struct {
	// DB
	DB_HOST     string `env:"DB_HOST" envDefault:"localhost"`
	DB_PORT     int    `env:"DB_PORT" envDefault:"5432"`
	DB_USER     string `env:"DB_USER" envDefault:"postgres"`
	DB_PASSWORD string `env:"DB_PASSWORD" secured:"true"`
	DB_NAME     string `env:"DB_NAME" envDefault:"default"`
	DB_SSL      string `env:"DB_SSL" envDefault:"disable"`
	// APP
	LOG_LEVEL  string `env:"LOG_LEVEL" envDefault:"debug"`
	LOG_FORMAT string `env:"LOG_FORMAT" envDefault:"plain"`
	VERSION    string `env:"VERSION" envDefault:"0.0.0"`
	APP_NAME   string `env:"APP_NAME" envDefault:"todo"`

	// FEATURES
	METRICS bool `env:"METRICS" envDefault:"true"`

	// HTTP
	PORT int `env:"PORT" envDefault:"8080"`
}

func (o *Configs) GetSecuredFilds() []string {
	secureFields := []string{}

	e := reflect.ValueOf(o).Elem()

	for i := 0; i < e.NumField(); i++ {
		field := e.Type().Field(i)

		if getStructTag(field, "secured") != "" {
			secureFields = append(secureFields, field.Name)
		}
	}

	return secureFields
}

func (o *Configs) GetFieldsWithValues() map[string]string {
	result := make(map[string]string)

	secureFields := o.GetSecuredFilds()

	e := reflect.ValueOf(o).Elem()
	for i := 0; i < e.NumField(); i++ {
		field := e.Type().Field(i)

		if helpers.InArray(field.Name, secureFields) {
			value := e.Field(i).String()

			switch lvalue := len(value); {
			case lvalue == 0:
				result[field.Name] = ""
			//nolint:gomnd
			case lvalue > 10:
				result[field.Name] = value[:6] + "..."
			case lvalue > 0:
				result[field.Name] = "xxxxxxx"
			}
		} else {
			result[field.Name] = fmt.Sprintf("%v", e.Field(i).Interface())
		}
	}

	return result
}

func getStructTag(f reflect.StructField, tagName string) string {
	return f.Tag.Get(tagName)
}

func (o *Configs) loadFromEnv() *Configs {
	options := Configs{}
	if err := env.Parse(&options); err != nil {
		panic(err)
	}

	return &options
}
