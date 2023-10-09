package config

import (
	"errors"
	"fmt"

	"github.com/debeando/go-common/file"
	"github.com/debeando/go-common/net"

	"github.com/go-yaml/yaml"
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
		return fmt.Errorf("Imposible to parse config file: %s", err)
	}

	return c.SanityCheck()
}

// SanityCheck verify the minimum config settings and set default values to start.
func (c *Config) SanityCheck() error {
	if c.General.Interval < 3 {
		c.General.Interval = 10

		return errors.New("Use positive value, and minimun start from 3 seconds, using default 10 seconds.")
	}

	if len(c.General.Hostname) == 0 {
		c.General.Hostname = net.Hostname()

		return errors.New("general.hostname: Custom value is not set, using current.")
	}

	return nil
}
