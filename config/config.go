package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/swapbyt3s/zenit/common"

	"github.com/go-yaml/yaml"
)

var File Config

// Init does any initialization necessary for the module.
func init() {
	File = Config {
		Path: "/etc/zenit/zenit.yaml",
		IPAddress: common.IPAddress(),
	}
}

func (c *Config) Load() error {
	source, err := ioutil.ReadFile(c.Path)
	if err != nil {
		source, err = ioutil.ReadFile("zenit.yaml")
		if err != nil {
			return errors.New(fmt.Sprintf("Fail to read config file: %s or %s", c.Path, "./zenit.yaml"))
		}
	}

	source = []byte(os.ExpandEnv(string(source)))

	if err := yaml.Unmarshal(source, &c); err != nil {
		errors.New(fmt.Sprintf("Imposible to parse config file - %s", err))
	}

	return nil
}

// SanityCheck verify the minimum config settings and set default values to start.
func (c *Config) SanityCheck() string {
	if c.General.Interval < 3 {
		c.General.Interval = 30

		return "Use positive value, and minimun start from 3 seconds, using default 30 seconds."
	}

	if len(c.General.Hostname) == 0 {
		c.General.Hostname = common.Hostname()

		return "general.hostname: Custom value is not set, using current."
	}

	return ""
}
