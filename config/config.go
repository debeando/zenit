package config

import (
	"github.com/debeando/go-common/file"
	"github.com/debeando/go-common/log"
	"github.com/debeando/go-common/net"

	"gopkg.in/yaml.v3"
)

var this *Config

func GetInstance() *Config {
	if this == nil {
		this = &Config{}
	}
	return this
}

func (c *Config) Load() {
	source := file.ReadExpandEnvAsString(c.Path)

	if err := yaml.Unmarshal([]byte(source), &c); err != nil {
		log.ErrorWithFields("Config", log.Fields{"message": "Imposible to parse config file", "file": c.Path})
		return
	}

	c.SanityCheck()
}

// SanityCheck verify the minimum config settings and set default values to start.
func (c *Config) SanityCheck() {
	if c.General.Interval < 3 {
		c.General.Interval = 10

		log.Error("Use positive value, and minimun start from 3 seconds, using default 10 seconds.")
	}

	if len(c.General.Hostname) == 0 {
		c.General.Hostname = net.Hostname()

		log.Error("general.hostname: Custom value is not set, using current.")
	}
}
