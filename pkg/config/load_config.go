package config

import (
	"fmt"
	"github.com/cristalhq/aconfig"
	"github.com/cristalhq/aconfig/aconfigdotenv"
	"os"
	"path"
	"reflect"
)

type Config interface {
	Validate() error
}

func LoadConfig(cfg Config) error {
	if reflect.ValueOf(cfg).Kind() != reflect.Ptr {
		return fmt.Errorf("config variable must be a pointer")
	}

	pwdDir, err := os.Getwd()
	if err != nil {
		return err
	}

	fileName := "dev.env"

	aconf := aconfig.Config{
		AllowUnknownFields: true,
		SkipFlags:          true,
		SkipDefaults:       false,
		Files:              []string{path.Join(pwdDir, fileName)},
		FileDecoders: map[string]aconfig.FileDecoder{
			".env": aconfigdotenv.New(),
		},
	}

	loader := aconfig.LoaderFor(cfg, aconf)
	if err := loader.Load(); err != nil {
		return err
	}

	return cfg.Validate()
}
