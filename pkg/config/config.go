// Package config global configuration of the program
package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Configuration struct {
	DatabaseFile string `required:"true" envconfig:"DATABASE_FILE" default:"/data/sqlite.db"`
	Source       string `required:"true" envconfig:"SOURCE_FILE" default:"/data/source.json"`
	AI_URL       string `required:"true" envconfig:"AI_URL"`
	AI_Model     string `required:"true" envconfig:"AI_MODEL"`
}

var config *Configuration = nil

func GetConfig() (*Configuration, error) {
	if config != nil {
		return config, nil
	}

	var s Configuration
	err := envconfig.Process("", &s)
	if err != nil {
		return nil, err
	}

	config = &s
	return config, nil
}
