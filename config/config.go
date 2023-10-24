package config

import (
	"github.com/debeando/go-common/file"
	"github.com/debeando/go-common/log"

	"gopkg.in/yaml.v3"
)

var this *Config

func GetInstance() *Config {
	if this == nil {
		this = &Config{}
	}
	return this
}

func (c *Config) Load() error {
	source := file.ReadExpandEnvAsString(c.Path)

	if err := yaml.Unmarshal([]byte(source), &c); err != nil {
		log.ErrorWithFields("Config", log.Fields{"error": err, "file": c.Path})
		return err
	}

	return nil
}
