package config

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/swapbyt3s/zenit/common"

	"github.com/sirupsen/logrus"
	"github.com/go-yaml/yaml"
)

var File Config

// Init does any initialization necessary for the module.
func init() {
	File = Config {
		Path: "/etc/zenit/zenit.yaml",
		IPAddress: common.IPAddress(),
	}

	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logrus.SetLevel(logrus.DebugLevel)
}

func (c *Config) Load() {
	source, err := ioutil.ReadFile(c.Path)
	if err != nil {
		source, err = ioutil.ReadFile("zenit.yaml")
		if err != nil {
			logrus.WithFields(map[string]interface{}{
				"error": fmt.Sprintf("Fail to read config file: %s or %s", c.Path, "./zenit.yaml"),
			}).Error("Config")
			os.Exit(1)
		}
	}

	source = []byte(os.ExpandEnv(string(source)))

	err = yaml.Unmarshal(source, &c)
	if err != nil {
		logrus.WithFields(map[string]interface{}{
			"error": fmt.Sprintf("Imposible to parse config file - %s", err),
		}).Error("Config")
		os.Exit(1)
	}
}

// SanityCheck verify the minimum config settings and set default values to start.
func (c *Config) SanityCheck() {
	if c.General.Interval < 3 {
		logrus.Warning("Config", map[string]interface{}{
			"message": "Use positive value, and minimun start from 3 seconds, using default 30 seconds.",
		})
		c.General.Interval = 30
	}

	if len(c.General.Hostname) == 0 {
		logrus.Warning("Config", map[string]interface{}{
			"message": "general.hostname: Custom value is not set, using current.",
		})

		c.General.Hostname = common.Hostname()
	}
}
